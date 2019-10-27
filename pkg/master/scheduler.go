package master

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	pb "gitlab.com/CBCTF/bullseye-runner/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type JobQ struct {
	job  *Job
	done chan struct{}
}

var (
	MasterCtx context.Context
	CancelMgr *CancelManager
	connPool  *ConnPool
	jqCh      chan JobQ
)

const (
	X11CapPrefix = "/var/tmp"
)

func sendCallback(url string, results []Result) error {
	if url == "" {
		return nil
	}
	buf, err := json.Marshal(results)
	if err != nil {
		return err
	}
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(buf))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode/100 != 2 { // != 2xx
		return fmt.Errorf("status is %d (!= 2xx)", resp.StatusCode)
	}

	return nil
}

// updateResult
func updateResult(db *gorm.DB) {
	resultMap := make(map[uint][]JobQ)
	callbackMap := make(map[string][]Result)

	cnt := len(jqCh)
	for cnt > 0 {
		jq := <-jqCh
		job := jq.job
		resultMap[job.ResultID] = append(resultMap[job.ResultID], jq)
		cnt--
	}

	for resultID, jqs := range resultMap {
		result := Result{}
		db.Find(&result, "id = ?", resultID)
		jobs := []Job{}
		for _, jq := range jqs {
			if jq.job.Succeeded {
				result.Succeeded++
			}
			result.Executed++
			jobs = append(jobs, *jq.job)
		}
		db.Save(&result)
		db.Model(&result).Association("Jobs").Append(jobs)
		for _, jq := range jqs {
			close(jq.done)
		}

		// send callback
		round := Round{}
		db.Find(&round, "id = ?", result.RoundID)
		url := round.CallbackURL
		if url != "" {
			callbackMap[url] = append(callbackMap[url], result)
		}
	}

	go func() {
		for url, results := range callbackMap {
			err := sendCallback(url, results)
			if err != nil {
				logger.Warn("callback", zap.Error(err))
			}
		}
	}()
}

func RunScheduler(db *gorm.DB) {
	// initialize
	MasterCtx = context.Background()
	CancelMgr = NewCancelManager()
	connPool = NewConnPool()
	jqCh = make(chan JobQ, 100000)
	rand.Seed(time.Now().UnixNano())

	ticker := time.NewTicker(time.Second * 10)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			err := doSchedule(db)
			if err != nil {
				logger.Warn("scheduling error", zap.Error(err))
			}
			updateResult(db)
		}
	}
}

func doSchedule(db *gorm.DB) error {
	var rounds []Round

	logger.Debug("checking rounds")

	// find past unexecuted round or manually added ones
	db.Preload("Results").Where("start_at <= ?", time.Now()).Or("start_at = NULL").Find(&rounds)

	for _, round := range rounds {
		// skip if already executed
		if len(round.Results) > 0 {
			continue
		}

		digest, err := func() (string, error) {
			// return specified hash if exists
			if round.ImageHash != "" {
				return round.ImageHash, nil
			}
			// get latest hash
			image, err := findImage(db, round)
			if err != nil {
				logger.Debug("couldn't find appropriate image")
				return "", err
			}
			logger.Debug("found image", zap.String("image", fmt.Sprintf("%+v", image)))
			round.ImageHash = image.Digest
			db.Save(&round)
			return image.Digest, nil
		}()
		if err != nil {
			continue
		}

		go func() {
			if err := doRound(db, round, digest); err != nil {
				logger.Warn("doRound", zap.Error(err))
			}
		}()
	}

	return nil
}

func findImage(db *gorm.DB, round Round) (*Image, error) {
	image := Image{}
	hit := 0

	db.Where("team = ? and exploit_container = ?", round.Team, round.ExploitContainer).
		Where("created_at <= ?", round.StartAt).
		Order("created_at desc", true).
		First(&image).Count(&hit)

	if hit == 0 {
		return nil, fmt.Errorf("Image not found")
	}

	return &image, nil
}

func doRound(db *gorm.DB, round Round, digest string) error {
	yml, err := EscapedTemplate(round.Yml, map[string]string{
		"exploitHash":      "@" + digest,
		"team":             round.Team,
		"exploitContainer": round.ExploitContainer,
	})
	if err != nil {
		return err
	}

	logger.Info("scheduling round", zap.Int("roundID", int(round.ID)))

	result := Result{}
	db.Model(&round).Association("Results").Append(&result)

	ctx, err := CancelMgr.Add(fmt.Sprintf("%d", round.ID), MasterCtx)
	if err != nil {
		return err
	}

	type WorkerInfo struct {
		host string
		ncpu uint64
	}
	workerInfos := []WorkerInfo{}

	err = func() error {
		logger.Debug("pulling image")
		for _, workerhost := range strings.Split(round.WorkerHosts, ",") {
			req := &pb.RunnerRequest{
				Uuid:             NewUUID(),
				Yml:              yml,
				RegistryHost:     round.RegistryHost,
				RegistryUsername: round.RegistryUsername,
				RegistryPassword: round.RegistryPassword,
				FlagTemplate:     round.FlagTemplate,
				PullImage:        true,
			}
			grpcCli, err := connPool.GetConn(workerhost)
			if err != nil {
				logger.Warn("failed to get grpc connection", zap.Error(err))
				return err
			}

			_, err = SendRequest(grpcCli, req, ctx)
			if err != nil {
				logger.Warn("PullRequest", zap.Error(err))
				return err
			}

			// update CPU info per host
			res, err := pb.NewRunnerClient(grpcCli).Info(ctx, &pb.InfoRequest{})
			if err != nil {
				logger.Warn("InfoRequest", zap.Error(err))
				return err
			}

			workerInfos = append(workerInfos, WorkerInfo{host: workerhost, ncpu: res.Cpus})
		}
		return nil
	}()
	if err != nil {
		logger.Warn("failed to pull image", zap.Error(err))
		return err
	}

	// add weight according to num of CPUs
	workerHosts := []string{}
	for _, workerInfo := range workerInfos {
		for i := 0; i < int(workerInfo.ncpu); i++ {
			workerHosts = append(workerHosts, workerInfo.host)
		}
	}

	for i := 0; i < int(round.Ntrials); i++ {
		workerhost := workerHosts[i%len(workerHosts)]

		uuid := NewUUID()
		// avoid conflication
		for ok := CancelMgr.Has(uuid); ok; uuid = NewUUID() {
		}

		_yml, err := EscapedTemplate(yml, map[string]int{
			"trialNumber": i,
		})
		if err != nil {
			return fmt.Errorf("failed to generate yml: %+v", err)
		}

		req := &pb.RunnerRequest{
			Uuid:             uuid,
			Timeout:          uint64(round.Timeout),
			Yml:              _yml,
			RegistryHost:     round.RegistryHost,
			RegistryUsername: round.RegistryUsername,
			RegistryPassword: round.RegistryPassword,
			FlagTemplate:     round.FlagTemplate,
		}

		if round.X11required {
			req.X11Info = &pb.X11Info{
				Width:  1024,
				Height: 768,
				Depth:  24,
				CapExt: "mp4",
			}
		}

		_ctx, _ := CancelMgr.Add(uuid, ctx) // uuid was checked beforehand

		go func() {
			defer CancelMgr.Cancel(uuid)

			grpcCli, err := connPool.GetConn(workerhost)
			if err != nil {
				logger.Warn("failed to get grpc connection", zap.Error(err))
				return
			}

			res, err := SendRequest(grpcCli, req, _ctx)
			if err != nil {
				logger.Warn("SendRequest", zap.Error(err))
				return
			}

			job := &Job{
				UUID:      req.Uuid,
				Host:      workerhost,
				Succeeded: res.Succeeded,
				Output:    res.Output,
				ResultID:  result.ID,
			}

			if round.X11required {
				err := ioutil.WriteFile(fmt.Sprintf("%s/%s.%s", X11CapPrefix, req.Uuid, req.X11Info.CapExt), res.X11Cap, 0644)
				if err != nil {
					logger.Warn("failed to save X11 capture", zap.Error(err))
				}
			}

			jq := JobQ{
				job:  job,
				done: make(chan struct{}, 0),
			}

			jqCh <- jq // send result to update daemon
			<-jq.done  // wait response from update daemon
		}()
	}

	return nil
}

func SendRequest(grpcCli *grpc.ClientConn, req *pb.RunnerRequest, ctx context.Context) (*pb.RunnerResponse, error) {
	res, err := pb.NewRunnerClient(grpcCli).Run(ctx, req)
	if err != nil {
		logger.Warn("grpc error", zap.Error(err))
		return nil, err
	}
	logger.Debug("response", zap.String("host", grpcCli.Target()), zap.String("uuid", res.Uuid), zap.Bool("succeeded", res.Succeeded))

	return res, nil
}

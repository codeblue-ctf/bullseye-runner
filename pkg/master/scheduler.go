package master

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	pb "gitlab.com/CBCTF/bullseye-runner/proto"
)

type JobQ struct {
	job  *Job
	done chan struct{}
}

var (
	MasterCtx context.Context
	CancelMgr *CancelManager
	jqCh      chan JobQ
)

// updateResult
func updateResult(db *gorm.DB) {
	resultMap := make(map[uint][]JobQ)

	cnt := len(jqCh)
	for cnt > 0 {
		jq := <-jqCh
		job := jq.job
		// log.Printf("%+v", job)
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
			jobs = append(jobs, *jq.job)
		}
		db.Save(&result)
		db.Model(&result).Association("Jobs").Append(jobs)
		for _, jq := range jqs {
			close(jq.done)
		}
	}
}

func RunScheduler(db *gorm.DB) {
	// initialize
	MasterCtx = context.Background()
	CancelMgr = NewCancelManager()
	jqCh = make(chan JobQ, 100000)
	rand.Seed(time.Now().UnixNano())

	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			err := doSchedule(db)
			if err != nil {
				log.Printf("schedule error: %v", err)
			}
			updateResult(db)
		}
	}
}

func doSchedule(db *gorm.DB) error {
	var rounds []Round

	log.Printf("checking rounds")

	// find past unexecuted round or manually added ones
	db.Preload("Results").Where("start_at <= ?", time.Now()).Or("start_at = NULL").Find(&rounds)

	for _, round := range rounds {
		// skip if already executed
		if len(round.Results) > 0 {
			continue
		}

		digest, err := func() (string, error) {
			// return specified hash if exists
			if round.ExploitHash != "" {
				return round.ExploitHash, nil
			}
			// get latest hash
			image, err := findImage(db, round)
			if err != nil {
				return "", err
			}
			return image.Digest, nil
		}()
		if err != nil {
			continue
		}

		if err := doRound(db, round, digest); err != nil {
			log.Printf("%+v", err)
		}
	}

	return nil
}

func findImage(db *gorm.DB, round Round) (*Image, error) {
	image := Image{}
	hit := 0

	db.Where("team = ? and problem = ?", round.Team, round.Problem).
		Where("created_at <= ?", round.StartAt).
		Order("created_at").
		First(&image).Count(&hit)

	if hit == 0 {
		return nil, fmt.Errorf("Image not found")
	}

	return &image, nil
}

func doRound(db *gorm.DB, round Round, digest string) error {
	yml, err := EscapedTemplate(round.Yml, map[string]string{
		"exploitHash": "@" + digest,
		"team":        round.Team,
		"problem":     round.Problem,
	})
	if err != nil {
		return err
	}

	log.Printf("scheduling round: %d", round.ID)

	result := Result{}
	db.Model(&round).Association("Results").Append(&result)

	ctx, err := CancelMgr.Add(fmt.Sprintf("%d", round.ID), MasterCtx)

	workerHosts := strings.Split(round.WorkerHosts, ",")
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

		_ctx, _ := CancelMgr.Add(uuid, ctx) // uuid was checked beforehand

		go func() {
			defer CancelMgr.Cancel(uuid)

			grpcCli, err := CreateGrpcCli(workerhost)
			defer grpcCli.Close()
			if err != nil {
				log.Printf("failed to create grpc connection: %+v", err)
				return
			}

			res, err := SendRequest(pb.NewRunnerClient(grpcCli), req, _ctx)
			if err != nil {
				log.Printf("error in SendRequest: %+v", err)
				return
			}

			job := &Job{
				UUID:      req.Uuid,
				Host:      workerhost,
				Succeeded: res.Succeeded,
				Output:    res.Output,
				ResultID:  result.ID,
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

func SendRequest(client pb.RunnerClient, req *pb.RunnerRequest, ctx context.Context) (*pb.RunnerResponse, error) {
	res, err := client.Run(ctx, req)
	if err != nil {
		log.Printf("%v.Run(_) = _, %v", client, err)
	}
	log.Printf("%+v", res)

	return res, nil
}

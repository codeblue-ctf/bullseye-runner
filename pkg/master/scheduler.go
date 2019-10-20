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

var (
	MasterCtx context.Context
	CancelMgr *CancelManager
	jobCh     chan Job
)

// updateResult
func updateResult(db *gorm.DB) {
	resultMap := make(map[uint][]Job)

	cnt := len(jobCh)
	log.Printf("updating: %d", len(jobCh))
	for cnt > 0 {
		job := <-jobCh
		log.Printf("%+v", job)
		resultMap[job.ResultID] = append(resultMap[job.ResultID], job)
		cnt--
	}

	for resultID, jobs := range resultMap {
		result := Result{}
		db.Where("id = ?", resultID).Find(&result)
		for _, job := range jobs {
			if job.Succeeded {
				result.Succeeded++
			}
		}
		db.Save(&result)
	}
}

func RunScheduler(db *gorm.DB) {
	// initialize
	MasterCtx = context.Background()
	CancelMgr = NewCancelManager()
	jobCh = make(chan Job, 100000)
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
			// updateResult(db)
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

	db.Where("team_id = ? and problem_id = ?", round.TeamID, round.ProblemID).
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

			job := &Job{
				UUID: req.Uuid,
				Done: false,
				Host: workerhost,
			}
			db.Model(&result).Association("Jobs").Append(&job)

			res, err := SendRequest(pb.NewRunnerClient(grpcCli), req, _ctx)
			if err != nil {
				log.Printf("error in SendRequest: %+v", err)
				return
			}

			job.Succeeded = res.Succeeded
			job.Output = res.Output
			job.Done = true

			db.Save(&job)

			// jobCh <- job
			log.Printf("added job: %+v", job)
		}()
	}

	return nil
}

func SendRequest(client pb.RunnerClient, req *pb.RunnerRequest, ctx context.Context) (*pb.RunnerResponse, error) {
	res, err := client.Run(ctx, req)
	if err != nil {
		log.Fatalf("%v.Run(_) = _, %v", client, err)
	}
	log.Printf("%+v", res)

	return res, nil
}

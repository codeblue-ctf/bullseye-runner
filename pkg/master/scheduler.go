package master

import (
	"context"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	pb "gitlab.com/CBCTF/bullseye-runner/proto"
)

var (
	mutex     sync.Mutex
	processes map[string]context.CancelFunc
)

func RunScheduler(db *gorm.DB) {
	processes = make(map[string]context.CancelFunc)

	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			err := doSchedule(db)
			if err != nil {
				log.Printf("schedule error: %v", err)
			}
		}
	}
}

// func cacheSchedule(db *gorm.DB) error {
// 	var schedules []Schedule
// 	db.Find(&schedules)

// 	processes = nil
// 	for _, sched := range schedules {
// 		now := time.Now()
// 		if !sched.Enabled {
// 			continue
// 		}
// 		if sched.StartAt.Before(now) || sched.StopAt.After(now) {
// 			continue
// 		}
// 		var count uint
// 		db.Where(Result{Schedule: sched}).Count(&count)
// 		if count == 0 {
// 			timer := time.NewTimer(time.Second * 10)
// 			defer timer.Stop()
// 			ctx, cancel := context.WithCancel(context.Background())

// 			go func() {
// 				for {
// 					select {
// 					case <-timer.C:
// 						log.Printf("executing sched: %+v", sched)
// 					case <-ctx.Done(): // canceled
// 						log.Printf("canceled schedule: %+v", sched)
// 						return
// 					}

// 				}
// 			}()
// 			process := Process{
// 				sched,
// 				cancel,
// 			}
// 			processes = append(processes, process)
// 		}
// 	}

// 	return nil
// }

func doSchedule(db *gorm.DB) error {
	var rounds []Round

	log.Printf("checking rounds")

	// find past unexecuted round
	db.Where("start_at <= ?", time.Now()).Find(&rounds)

	for _, round := range rounds {
		// skip if already executed
		results := []Result{}
		db.Model(&round).Related(&results)
		if len(results) > 0 {
			continue
		}

		if err := doRound(db, round); err != nil {
			log.Printf("%+v", err)
		}
	}

	return nil
}

func doRound(db *gorm.DB, round Round) error {
	// get latest hash
	record := DockerHash{}
	hit := 0
	db.Where("team_id = ? and problem_id = ?", round.TeamID, round.ProblemID).Where("timestamp <= ?", round.StartAt).Order("timestamp").First(&record).Count(&hit)
	Debug(db.Where("team_id = ? and problem_id = ?", round.TeamID, round.ProblemID).Where("timestamp <= ?", round.StartAt).Order("timestamp").First(&record).Count(&hit).QueryExpr())
	if hit == 0 {
		return nil
	}

	log.Printf("found: %s\n", record.Digest)
	yml, err := EscapedTemplate(round.Yml, map[string]string{
		"exploitHash": record.Digest,
	})
	if err != nil {
		return err
	}

	log.Printf("scheduling round: %d\n", round.ID)

	result := Result{}
	db.Model(&round).Association("Results").Append(&result)

	workerHosts := strings.Split(round.WorkerHosts, ",")
	for i := 0; i < int(round.Ntrials); i++ {
		workerhost := workerHosts[i%len(workerHosts)]

		uuid, err := NewUUID()
		for _, ok := processes[uuid]; ok != true; uuid, _ = NewUUID() {
		}

		if err != nil {
			return err
		}

		req := pb.RunnerRequest{
			Uuid:          uuid,
			Timeout:       uint64(round.Timeout),
			Yml:           yml,
			RegistryToken: "test",
			FlagTemplate:  round.FlagTemplate,
		}

		ctx, cancel := context.WithCancel(context.Background())
		if _, ok := processes[uuid]; ok == true {
			log.Printf("worker %s is ongoing", uuid)
			continue
		}

		mutex.Lock()
		processes[uuid] = cancel
		mutex.Unlock()

		go func(ctx context.Context, host string, req *pb.RunnerRequest, result *Result) error {
			defer func() { // cleanup
				mutex.Lock()
				if _, ok := processes[uuid]; ok == true {
					delete(processes, uuid)
				}
				mutex.Unlock()
			}()

			grpcCli, err := CreateGrpcCli(workerhost)
			if err != nil {
				return err
			}

			res, err := SendRequest(pb.NewRunnerClient(grpcCli), req, ctx)
			if err != nil {
				log.Printf("%v.Run(_) = _, %+v", err)
			}
			grpcCli.Close()

			job := Job{
				UUID:      uuid,
				Succeeded: res.Succeeded,
				Output:    res.Output,
			}
			db.Model(&result).Association("Jobs").Append(&job)

			return nil
		}(ctx, workerhost, &req, &result)
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

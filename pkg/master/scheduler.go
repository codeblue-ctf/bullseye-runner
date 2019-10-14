package master

import (
	"context"
	"log"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	pb "gitlab.com/CBCTF/bullseye-runner/proto"
)

var (
	processes map[uint]context.CancelFunc
)

func RunScheduler(db *gorm.DB) {
	pool := &ConnPool{}
	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			err := doSchedule(db, pool)
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

func doSchedule(db *gorm.DB, pool *ConnPool) error {
	var rounds []Round

	// find past unexecuted round
	db.Select(rounds).Where("start_at <= ?", time.Now())

	for _, round := range rounds {
		// get latest hash
		record := &DockerHash{}
		db.Where("team_id == ? and problem_id == ?", round.TeamID, round.ProblemID).Order("timestamp").First(record)
		if record == nil {
			continue
		}

		yml, err := EscapedTemplate(round.Yml, map[string]string{
			"containerHash": record.Digest,
		})
		if err != nil {
			return err
		}

		for i := 0; i < int(round.Ntrials); i++ {
			workerhost := round.WorkerHosts[i%len(round.WorkerHosts)]

			pbcli, err := pool.GetConn(workerhost)
			if err != nil {
				return err
			}

			result := Result{
				Round: round,
			}
			db.Create(result)

			uuid := "hoge" // TODO

			req := &pb.RunnerRequest{
				Uuid:          uuid,
				Timeout:       uint64(round.Timeout),
				Yml:           yml,
				RegistryToken: "test",
				FlagTemplate:  round.FlagTemplate,
			}

			ctx, cancel := context.WithCancel(context.Background())
			if _, ok := processes[round.ID]; ok == true {
				log.Printf("round %d is ongoing", round.ID)
				continue
			}
			processes[round.ID] = cancel

			go func(ctx context.Context) {
				res, err := SendRequest(pb.NewRunnerClient(pbcli), req, ctx)
				if err != nil {
					log.Printf("%v.Run(_) = _, %+v", err)
				}
				workerResult := WorkerResult{
					Uuid:      uuid,
					Succeeded: res.Succeeded,
					Output:    res.Output,
					Result:    result,
				}
				db.Create(workerResult)
			}(ctx)
		}
	}

	return nil
}

func SendRequest(client pb.RunnerClient, req *pb.RunnerRequest, ctx context.Context) (*pb.RunnerResponse, error) {
	// TODO: replace mock
	log.Printf("%+v", req)
	return nil, nil

	res, err := client.Run(ctx, req)
	if err != nil {
		log.Fatalf("%v.Run(_) = _, %v", client, err)
	}
	log.Printf("%+v", res)

	return nil, nil
}

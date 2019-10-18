package master

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	pb "gitlab.com/CBCTF/bullseye-runner/proto"
)

var (
	masterCtx     context.Context
	cancelManager *CancelManager
)

type CancelManager struct {
	mut sync.Mutex
	c   map[string]context.CancelFunc
}

func NewCancelManager() *CancelManager {
	return &CancelManager{
		c: make(map[string]context.CancelFunc),
	}
}

func (cm *CancelManager) Has(key string) bool {
	cm.mut.Lock()
	defer cm.mut.Unlock()
	_, ok := cm.c[key]
	return ok
}

func (cm *CancelManager) Add(key string, _ctx context.Context) (context.Context, error) {
	cm.mut.Lock()
	defer cm.mut.Unlock()
	if _, ok := cm.c[key]; ok {
		return nil, fmt.Errorf("key %s already exists", key)
	}

	ctx, cancel := context.WithCancel(_ctx)
	cm.c[key] = cancel
	return ctx, nil
}

func (cm *CancelManager) Cancel(key string) error {
	cm.mut.Lock()
	defer cm.mut.Unlock()
	cancel, ok := cm.c[key]
	if !ok {
		return fmt.Errorf("key %s does not exist", key)
	}
	cancel()
	delete(cm.c, key)
	return nil
}

func RunScheduler(db *gorm.DB) {
	// initialize
	masterCtx = context.Background()
	cancelManager = NewCancelManager()
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
		}
	}
}

func doSchedule(db *gorm.DB) error {
	var rounds []Round

	log.Printf("checking rounds")

	// find past unexecuted round
	db.Preload("Results").Where("start_at <= ?", time.Now()).Find(&rounds)

	for _, round := range rounds {
		// skip if already executed
		if len(round.Results) > 0 {
			continue
		}

		// get latest hash
		image, err := findImage(db, round)
		if err != nil {
			continue
		}
		log.Printf("found: %s\n", image.Digest)

		if err := doRound(db, round, *image); err != nil {
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

func doRound(db *gorm.DB, round Round, image Image) error {
	yml, err := EscapedTemplate(round.Yml, map[string]string{
		"exploitHash": image.Digest,
	})
	if err != nil {
		return err
	}

	log.Printf("scheduling round: %d\n", round.ID)

	result := Result{}
	db.Model(&round).Association("Results").Append(&result)

	ctx, err := cancelManager.Add(fmt.Sprintf("%d"), masterCtx)

	workerHosts := strings.Split(round.WorkerHosts, ",")
	for i := 0; i < int(round.Ntrials); i++ {
		workerhost := workerHosts[i%len(workerHosts)]

		uuid := NewUUID()
		// avoid conflication
		for ok := cancelManager.Has(uuid); ok; uuid = NewUUID() {
		}

		req := &pb.RunnerRequest{
			Uuid:          uuid,
			Timeout:       uint64(round.Timeout),
			Yml:           yml,
			RegistryToken: "test",
			FlagTemplate:  round.FlagTemplate,
		}

		_ctx, _ := cancelManager.Add(uuid, ctx) // uuid was checked beforehand

		go func() {
			defer cancelManager.Cancel(uuid)

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

			job := Job{
				UUID:      req.Uuid,
				Succeeded: res.Succeeded,
				Output:    res.Output,
			}

			db.Model(&result).Association("Jobs").Append(&job)
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

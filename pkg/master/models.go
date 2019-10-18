package master

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Schedule struct {
	gorm.Model
	StartAt      time.Time `json:"start_at"`
	StopAt       time.Time `json:"stop_at"`
	Yml          string    `json:"yml"`
	FlagTemplate string    `json:"flag_template"`
	Interval     uint      `json:"interval"`
	Ntrials      uint      `json:"ntrials"`
	Timeout      uint      `json:"timeout"`
	WorkerHosts  string    `json:"worker_hosts"`
	CallbackURL  string    `json:"callback_url"`
	ProblemID    string    `json:"problem_id"`
	TeamID       string    `json:"team_id"`
	Rounds       []Round   `json:"rounds,omitempty"`
}

type Round struct {
	gorm.Model
	StartAt      time.Time `json:"start_at"`
	Yml          string    `json:"yml"`
	FlagTemplate string    `json:"flag_template"`
	Ntrials      uint      `json:"ntrials"`
	Timeout      uint      `json:"timeout"`
	WorkerHosts  string    `json:"worker_hosts"`
	CallbackURL  string    `json:"callback_url"`
	ProblemID    string    `json:"problem_id"`
	TeamID       string    `json:"team_id"`
	ScheduleID   uint      `json:"schedule_id"`
	Results      []Result  `json:"-,omitempty"`
}

type Result struct {
	gorm.Model
	Succeeded uint  `json:"succeeded"`
	RoundID   uint  `json:"round_id"`
	Jobs      []Job `json:"jobs"`
}

type Job struct {
	gorm.Model
	UUID      string `json:"uuid"`
	Host      string `json:"host"`
	Succeeded bool   `json:"succeeded"`
	Output    string `json:"output"`
	ResultID  uint   `json:"result_id"`
}

type DockerHash struct {
	gorm.Model
	UUID       string    `json:"uuid"`
	Timestamp  time.Time `json:"timestamp"`
	Digest     string    `json:"digest"`
	TeamID     string    `json:"team_id"`
	ProblemID  string    `json:"problem_id"`
	RemoteAddr string    `json:"remote_addr"`
	UserAgent  string    `json:"user_agent"`
}

func (s *Schedule) AfterCreate(db *gorm.DB) error {
	// create rounds according to schedule
	rounds := []Round{}
	for t := s.StartAt; t.Before(s.StopAt); t = t.Add(time.Duration(s.Interval) * time.Minute) {
		round := Round{
			StartAt:      t,
			Yml:          s.Yml,
			FlagTemplate: s.FlagTemplate,
			Ntrials:      s.Ntrials,
			Timeout:      s.Timeout,
			WorkerHosts:  s.WorkerHosts,
			CallbackURL:  s.CallbackURL,
			ProblemID:    s.ProblemID,
			TeamID:       s.TeamID,
		}
		rounds = append(rounds, round)
	}
	db.Model(s).Association("Rounds").Append(rounds)
	return nil
}

func (s *Schedule) BeforeDelete(db *gorm.DB) error {
	// cleanup unexecuted rounds
	rounds := []Round{}
	db.Model(s).Association("Rounds").Find(&rounds)
	db.Delete(&rounds)
	return nil
}

package master

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
)

type Schedule struct {
	gorm.Model
	StartAt          time.Time `json:"start_at"`
	StopAt           time.Time `json:"stop_at"`
	Yml              string    `json:"yml"`
	FlagTemplate     string    `json:"flag_template"`
	Interval         uint      `json:"interval"`
	Ntrials          uint      `json:"ntrials"`
	Timeout          uint      `json:"timeout"`
	WorkerHosts      string    `json:"worker_hosts"`
	RegistryHost     string    `json:"registry_host"`
	RegistryUsername string    `json:"registry_username"`
	RegistryPassword string    `json:"registry_password"`
	CallbackURL      string    `json:"callback_url"`
	ProblemID        string    `json:"problem_id"`
	TeamID           string    `json:"team_id"`
	Rounds           []Round   `json:"rounds,omitempty"`
}

type Round struct {
	gorm.Model
	StartAt          *time.Time `json:"start_at"`
	Yml              string     `json:"yml"`
	FlagTemplate     string     `json:"flag_template"`
	Ntrials          uint       `json:"ntrials"`
	Timeout          uint       `json:"timeout"`
	WorkerHosts      string     `json:"worker_hosts"`
	RegistryHost     string     `json:"registry_host"`
	RegistryUsername string     `json:"registry_username"`
	RegistryPassword string     `json:"registry_password"`
	CallbackURL      string     `json:"callback_url"`
	ProblemID        string     `json:"problem_id"`
	TeamID           string     `json:"team_id"`
	ScheduleID       uint       `json:"schedule_id"`
	ExploitHash      string     `json:"exploit_hash,omitempty"`
	Results          []Result   `json:"-"`
}

type Result struct {
	gorm.Model
	Succeeded uint  `json:"succeeded"`
	RoundID   uint  `json:"round_id"`
	Jobs      []Job `json:"-"`
}

type Job struct {
	gorm.Model
	UUID      string `json:"uuid"`
	Done      bool   `json:"done"`
	Host      string `json:"host"`
	Succeeded bool   `json:"succeeded"`
	Output    string `json:"output"`
	ResultID  uint   `json:"result_id"`
}

type Image struct {
	gorm.Model
	UUID       string `json:"uuid"`
	Digest     string `json:"digest"`
	TeamID     string `json:"team_id"`
	ProblemID  string `json:"problem_id"`
	RemoteAddr string `json:"remote_addr"`
	UserAgent  string `json:"user_agent"`
}

func (s *Schedule) AfterCreate(db *gorm.DB) error {
	// create rounds according to schedule
	rounds := []Round{}
	for t := s.StartAt; t.Before(s.StopAt); t = t.Add(time.Duration(s.Interval) * time.Minute) {
		_t := t
		round := Round{
			StartAt:          &_t,
			Yml:              s.Yml,
			FlagTemplate:     s.FlagTemplate,
			Ntrials:          s.Ntrials,
			Timeout:          s.Timeout,
			WorkerHosts:      s.WorkerHosts,
			RegistryHost:     s.RegistryHost,
			RegistryUsername: s.RegistryUsername,
			RegistryPassword: s.RegistryPassword,
			CallbackURL:      s.CallbackURL,
			ProblemID:        s.ProblemID,
			TeamID:           s.TeamID,
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

func (r *Result) BeforeDelete(db *gorm.DB) error {
	CancelMgr.Cancel(fmt.Sprintf("%d", r.ID))
	jobs := []Job{}
	db.Model(r).Association("Jobs").Find(&jobs)
	db.Delete(&jobs)
	return nil
}

func (j *Job) BeforeDelete(db *gorm.DB) error {
	CancelMgr.Cancel(j.UUID)
	return nil
}

package master

import (
	"fmt"
	"sync"
	"time"

	"github.com/jinzhu/gorm"
)

type Schedule struct {
	gorm.Model
	StartAt          time.Time `json:"start_at"`
	StopAt           time.Time `json:"stop_at"`
	Yml              string    `json:"yml" sql:"type:text;"`
	X11required      bool      `json:"x11_required"`
	FlagTemplate     string    `json:"flag_template"`
	Interval         uint      `json:"interval"`
	Ntrials          uint      `json:"ntrials"`
	Timeout          uint      `json:"timeout"`
	WorkerHosts      string    `json:"worker_hosts"`
	RegistryHost     string    `json:"registry_host"`
	RegistryUsername string    `json:"registry_username"`
	RegistryPassword string    `json:"registry_password"`
	CallbackURL      string    `json:"callback_url"`
	ExploitContainer string    `json:"exploit_container"`
	Team             string    `json:"team"`
	Rounds           []Round   `json:"rounds,omitempty"`
}

type Round struct {
	gorm.Model
	StartAt          *time.Time `json:"start_at"`
	Yml              string     `json:"yml" sql:"type:text;"`
	X11required      bool       `json:"x11_required"`
	FlagTemplate     string     `json:"flag_template"`
	Ntrials          uint       `json:"ntrials"`
	Timeout          uint       `json:"timeout"`
	WorkerHosts      string     `json:"worker_hosts"`
	RegistryHost     string     `json:"registry_host"`
	RegistryUsername string     `json:"registry_username"`
	RegistryPassword string     `json:"registry_password"`
	CallbackURL      string     `json:"callback_url"`
	ExploitContainer string     `json:"exploit_container"`
	Team             string     `json:"team"`
	ScheduleID       uint       `json:"schedule_id"`
	ImageHash        string     `json:"image_hash"`
	Checked          bool       `json:"checked"`
	Results          []Result   `json:"-"`
}

type Result struct {
	gorm.Model
	Succeeded uint  `json:"succeeded"`
	Executed  uint  `json:"executed"`
	RoundID   uint  `json:"round_id"`
	Jobs      []Job `json:"-"`
}

type Job struct {
	gorm.Model
	UUID      string `json:"uuid"`
	Host      string `json:"host"`
	Succeeded bool   `json:"succeeded"`
	Output    string `json:"output"`
	ResultID  uint   `json:"result_id"`
}

type Image struct {
	gorm.Model
	UUID             string `json:"uuid"`
	Digest           string `json:"digest"`
	Team             string `json:"team"`
	ExploitContainer string `json:"exploit_container"`
	RemoteAddr       string `json:"remote_addr"`
	UserAgent        string `json:"user_agent"`
}

var smut sync.Mutex

func (s *Schedule) AfterCreate(db *gorm.DB) error {
	// create rounds according to schedule
	rounds := []Round{}
	for t := s.StartAt; t.Before(s.StopAt); t = t.Add(time.Duration(s.Interval) * time.Minute) {
		_t := t
		round := Round{
			StartAt:          &_t,
			Yml:              s.Yml,
			X11required:      s.X11required,
			FlagTemplate:     s.FlagTemplate,
			Ntrials:          s.Ntrials,
			Timeout:          s.Timeout,
			WorkerHosts:      s.WorkerHosts,
			RegistryHost:     s.RegistryHost,
			RegistryUsername: s.RegistryUsername,
			RegistryPassword: s.RegistryPassword,
			CallbackURL:      s.CallbackURL,
			ExploitContainer: s.ExploitContainer,
			Team:             s.Team,
		}
		rounds = append(rounds, round)
	}
	db.Model(s).Association("Rounds").Append(rounds)
	return nil
}

func (s *Schedule) AfterDelete(db *gorm.DB) error {
	// cleanup unexecuted rounds
	// rounds := []Round{}
	// db.Model(s).Association("Rounds").Find(&rounds)
	// db.Delete(&rounds)
	return nil
}

func (r *Round) AfterDelete(db *gorm.DB) error {
	CancelMgr.Cancel(fmt.Sprintf("%d", r.ID))
	results := []Result{}
	db.Model(r).Related(&results)
	for _, result := range results {
		if result.ID == 0 {
			continue
		}
		db.Delete(&result)
	}
	return nil
}

func (r *Result) AfterDelete(db *gorm.DB) error {
	jobs := []Job{}
	db.Model(r).Association("Jobs").Find(&jobs)
	db.Delete(&jobs)
	CancelMgr.Cancel(fmt.Sprintf("%d", r.RoundID))
	return nil
}

func (j *Job) BeforeDelete(db *gorm.DB) error {
	CancelMgr.Cancel(j.UUID)
	return nil
}

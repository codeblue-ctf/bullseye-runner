package master

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Schedule struct {
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
	Enabled      bool      `json:"enabled"`
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
}

type Result struct {
	gorm.Model
	Succeeded  uint
	Round      Round `gorm:"foreignkey:RoundRefer"`
	RoundRefer uint
}

type WorkerResult struct {
	gorm.Model
	Uuid        string
	Succeeded   bool
	Output      string
	Result      Result `gorm:"foreignkey:ResultRefer"`
	ResultRefer uint
}

type DockerHash struct {
	gorm.Model
	Id         string    `json:"id"`
	Timestamp  time.Time `json:"timestamp"`
	Digest     string    `json:"digest"`
	TeamID     string    `json:"team_id"`
	ProblemID  string    `json:"problem_id"`
	RemoteAddr string    `json:"remote_addr"`
	UserAgent  string    `json:"user_agent"`
}

package master

import (
	"github.com/jinzhu/gorm"
)

type Schedule struct {
	gorm.Model
	StartAt      string   `json:"start_at"`
	StopAt       string   `json:"stop_at"`
	Yml          string   `json:"yml"`
	FlagTemplate string   `json:"flag_template"`
	Interval     uint     `json:"interval"`
	Ntrials      uint     `json:"ntrials"`
	WorkerHosts  []string `json:"worker_hosts"`
	CallbackUrl  string   `json:"callback_url"`
	Enabled      bool
}

type Result struct {
	gorm.Model
	succeeded     uint
	failed        uint
	output        string
	Schedule      Schedule `gorm:"foreignkey:ScheduleRefer"`
	ScheduleRefer uint
}

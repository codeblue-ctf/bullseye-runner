package master

type Schedule struct {
	CallBack string
	StartAt string
	StopAt string
	DockerComposeYml string
	Interval int
	Ntrials int
	WorkerHost string
}

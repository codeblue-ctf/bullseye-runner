package main

import (
	"bytes"
	"context"
	"io/ioutil"
	"log"
	"os/exec"
	"path"
	"strings"
	"time"

	pb "gitlab.com/CBCTF/bullseye-runner/proto"
)

func setupDirectory(yml string, flag string) (string, error) {
	dir, err := ioutil.TempDir("tmp", "bullseye-runner-")
	if err != nil {
		return "", err
	}

	err = ioutil.WriteFile(path.Join(dir, "docker-compose.yml"), []byte(yml), 0644)
	if err != nil {
		log.Printf("failed to write yml under %s", dir)
		return "", err
	}

	err = ioutil.WriteFile(path.Join(dir, "flag"), []byte(flag), 0644)
	if err != nil {
		log.Printf("failed to write flag under %s", dir)
		return "", err
	}

	return dir, nil
}

func runDockerCompose(dir string, timeout int32) (bool, string, error) {
	log.Printf("start evaluation: %s", dir)

	exec.Command("docker-compose", "up", "-d").Run()
	time.Sleep(time.Duration(timeout) * time.Millisecond)
	exec.Command("docker-compose", "kill").Run()
	cmd := exec.Command("docker-compose", "logs")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Printf("failed to get log")
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(stdout)

	exec.Command("docker-compose", "down").Run()

	success, err := checkFlag(dir)
	if err != nil {
		return false, "", err
	}

	log.Printf("end evaluation: %s", dir)

	if success {
		return true, buf.String(), nil
	}

	return false, buf.String(), nil
}

func checkFlag(dir string) (bool, error) {
	submittedFlag, err := ioutil.ReadFile(path.Join(dir, "submitted-flag"))
	if err != nil {
		log.Printf("failed to open submitted flag")
		return false, err
	}

	flag, err := ioutil.ReadFile(path.Join(dir, "flag"))
	if err != nil {
		log.Printf("failed to open flag")
		return false, err
	}

	trimmed := strings.TrimSpace(string(submittedFlag))
	trimmed = strings.Trim(trimmed, "\r\n")

	if trimmed == string(flag) {
		return true, nil
	}

	return false, nil
}

// run docker-compose.yml and return results as RunnerResponse
func RunRequest(ctx context.Context, req *pb.RunnerRequest) (*pb.RunnerResponse, error) {
	dir, err := setupDirectory(req.DockerComposeYml, req.FlagTemplate)
	if err != nil {
		return nil, err
	}

	succeeded, output, err := runDockerCompose(dir, req.Timeout)
	if err != nil {
		log.Printf("error: %v", err)
		return nil, err
	}

	res := pb.RunnerResponse{
		Uuid:      req.Uuid,
		Succeeded: succeeded,
		Output:    output,
	}

	return &res, nil
}

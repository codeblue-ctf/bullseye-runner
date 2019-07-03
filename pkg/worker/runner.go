package worker

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
	dir, err := ioutil.TempDir("./tmp", "bullseye-runner-")
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

func runCommand(args ...string) error {
	cmd := exec.Command(args[0], args[1:]...)
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	err := cmd.Run()

	log.Printf("%s %s, stdout: %#v", args[0], args[1:], stdout.String())
	log.Printf("%s %s, stderr: %#v", args[0], args[1:], stderr.String())

	if err != nil {
		return err
	}

	return nil
}

func runDockerCompose(dir string, timeout int32) (bool, string, error) {
	log.Printf("start evaluation: %s", dir)

	err := runCommand("docker-compose", "up", "-d")
	if err != nil {
		return false, "", err
	}

	time.Sleep(time.Duration(timeout) * time.Millisecond)

	err = runCommand("docker-compose", "kill")
	if err != nil {
		return false, "", err
	}

	cmd := exec.Command("docker-compose", "logs")
	output := new(bytes.Buffer)
	cmd.Stdout = output

	err = cmd.Run()
	if err != nil {
		return false, "", err
	}

	err = runCommand("docker-compose", "down")
	if err != nil {
		return false, "", err
	}

	success, err := checkFlag(dir)
	if err != nil {
		return false, "", err
	}

	log.Printf("end evaluation: %s", dir)

	if success {
		return true, output.String(), nil
	}

	return false, output.String(), nil
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

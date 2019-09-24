package worker

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"

	"github.com/lucasjones/reggen"
	pb "gitlab.com/CBCTF/bullseye-runner/proto"
)

const (
	TempDir      = "/tmp"
	FlagSuffix   = "flag"
	SubmitSuffix = "submit"
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

	curdir, err := os.Getwd()
	if err != nil {
		return false, "", err
	}
	if err := os.Chdir(dir); err != nil {
		return false, "", err
	}
	defer os.Chdir(curdir)

	err = runCommand("docker-compose", "up", "-d")
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

	// success, err := checkFlag(dir)
	// if err != nil {
	// 	return false, "", err
	// }

	// log.Printf("end evaluation: %s", dir)

	// if success {
	// 	return true, output.String(), nil
	// }

	return false, output.String(), nil
}

func GenerateFlag(template string) (string, error) {
	g, err := reggen.NewGenerator(template)
	if err != nil {
		return "", err
	}
	return g.Generate(10), nil
}

func trim(s []byte) string {
	return strings.Trim(fmt.Sprintf("%s", s), " \x00\r\n")
}

func GetFlagPaths(uuid string) (string, string) {
	return fmt.Sprintf("%s/%s-%s", TempDir, uuid, FlagSuffix), fmt.Sprintf("%s/%s-%s", TempDir, uuid, SubmitSuffix)
}

func CheckFlag(uuid string) (bool, error) {
	flagPath, submitPath := GetFlagPaths(uuid)

	submitBytes, err := ioutil.ReadFile(submitPath)
	if err != nil {
		log.Printf("failed to open submitted flag")
		return false, err
	}

	flagBytes, err := ioutil.ReadFile(flagPath)
	if err != nil {
		log.Printf("failed to open flag")
		return false, err
	}
	submitStr := trim(submitBytes)
	flagStr := trim(flagBytes)

	if submitStr == flagStr {
		log.Printf("correct flag: %s", flagStr)
		return true, nil
	}
	log.Printf("incorrect flag: %#v (!= %#v)", submitStr, flagStr)

	return false, nil
}

func Cleanup(uuid string) error {
	flagPath, submitPath := GetFlagPaths(uuid)
	if err := os.Remove(flagPath); err != nil {
		return err
	}
	if err := os.Remove(submitPath); err != nil {
		return err
	}
	return nil
}

// run received yml and return results as RunnerResponse
func RunRequest(ctx context.Context, req *pb.RunnerRequest) (*pb.RunnerResponse, error) {
	succeeded, output, err := RunDockerCompose(ctx, req)
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

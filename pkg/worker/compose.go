package worker

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"text/template"
	"time"

	"github.com/docker/cli/cli/config/configfile"
	clitypes "github.com/docker/cli/cli/config/types"
	apitypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/libcompose/docker"
	"github.com/docker/libcompose/docker/auth"
	dockerctx "github.com/docker/libcompose/docker/ctx"
	"github.com/docker/libcompose/project"
	"github.com/docker/libcompose/project/options"

	pb "gitlab.com/CBCTF/bullseye-runner/proto"
)

func RunDockerCompose(ctx context.Context, req *pb.RunnerRequest) (bool, string, error) {
	flagPath, submitPath := GetFlagPaths(req.Uuid)

	// generate flag from regex
	flagStr, err := GenerateFlag(req.FlagTemplate)
	if err != nil {
		return false, "", err
	}
	log.Printf("generated flag: %s", flagStr)

	flagFile, err := os.Create(flagPath)
	if err != nil {
		return false, "", err
	}
	_, err = flagFile.WriteString(flagStr)
	if err != nil {
		return false, "", err
	}
	if err = flagFile.Close(); err != nil {
		return false, "", err
	}

	submitFile, err := os.Create(submitPath)
	if err = submitFile.Close(); err != nil {
		return false, "", err
	}

	defer Cleanup(req.Uuid)

	var yml bytes.Buffer
	tpl, err := template.New("yml").Parse(req.DockerComposeYml)
	if err != nil {
		return false, "", err
	}

	dict := map[string]string{
		"flagPath":   flagPath,
		"submitPath": submitPath,
	}

	err = tpl.Execute(&yml, dict)
	if err != nil {
		return false, "", err
	}

	log.Printf("execute yml: %#v", yml.String())

	configFile := &configfile.ConfigFile{
		AuthConfigs: map[string]clitypes.AuthConfig{
			"localhost:5000": clitypes.AuthConfig{
				Username: "admin",
				Password: "password",
				// Auth:          req.DockerRegistryToken,
				ServerAddress: "localhost:5000",
			},
		},
	}

	project, err := docker.NewProject(&dockerctx.Context{
		Context: project.Context{
			ComposeBytes: [][]byte{yml.Bytes()},
			ProjectName:  req.Uuid,
		},
		ConfigFile: configFile,
		AuthLookup: auth.NewConfigLookup(configFile),
	}, nil)
	if err != nil {
		return false, "", err
	}

	// create network in advance to make evaluation faster
	client, err := client.NewEnvClient()
	if err != nil {
		return false, "", err
	}
	networkName := fmt.Sprintf("%s_default", req.Uuid)
	_, err = client.NetworkCreate(ctx, networkName, apitypes.NetworkCreate{})
	defer client.NetworkRemove(ctx, networkName)
	if err != nil {
		return false, "", err
	}

	err = project.Up(ctx, options.Up{})
	if err != nil {
		log.Printf("failed to up: %v", err)
		return false, "", err
	}

	time.Sleep(time.Duration(req.Timeout) * time.Millisecond)

	err = project.Delete(ctx, options.Delete{
		RemoveVolume:  true,
		RemoveRunning: true,
	})
	if err != nil {
		return false, "", err
	}

	err = project.Kill(ctx, "KILL")
	if err != nil {
		return false, "", err
	}

	ok, err := CheckFlag(req.Uuid)

	if err != nil {
		return false, "", err
	}
	if ok {
		return true, "", nil
	}

	return false, "", nil
}

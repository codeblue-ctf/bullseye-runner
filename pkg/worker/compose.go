package worker

import (
	"bytes"
	"context"
	"fmt"
	"log"
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

func getNeworkID(uuid string) string {
	return fmt.Sprintf("%s_default", uuid)
}

func prepareNetwork(ctx context.Context, req *pb.RunnerRequest) error {
	networkID := getNeworkID(req.Uuid)
	client, err := client.NewEnvClient()
	if err != nil {
		log.Printf("failed to create env client: %v", err)
		return err
	}
	_, err = client.NetworkCreate(ctx, networkID, apitypes.NetworkCreate{})
	if err != nil {
		log.Printf("failed to create network: %v", err)
		return err
	}

	return nil
}

func cleanNetwork(req *pb.RunnerRequest) error {
	networkID := getNeworkID(req.Uuid)
	client, err := client.NewEnvClient()
	if err != nil {
		return err
	}
	if err := client.NetworkRemove(context.Background(), networkID); err != nil {
		return err
	}

	return nil
}

func cleanCompose(req *pb.RunnerRequest, project project.APIProject) {
	project.Delete(context.Background(), options.Delete{
		RemoveVolume:  true,
		RemoveRunning: true,
	})
	cleanNetwork(req)
}

func RunDockerCompose(ctx context.Context, req *pb.RunnerRequest) (bool, string, error) {
	flagPath, submitPath, err := PrepareFlags(req.Uuid, req.FlagTemplate)
	if err != nil {
		return false, "", err
	}
	defer CleanFlags(req.Uuid)

	var yml bytes.Buffer
	tpl, err := template.New("yml").Parse(req.Yml)
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
			req.RegistryHost: clitypes.AuthConfig{
				Username: req.RegistryUsername,
				Password: req.RegistryPassword,
				// Auth:          req.DockerRegistryToken,
				ServerAddress: req.RegistryHost,
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
	if err := prepareNetwork(ctx, req); err != nil {
		return false, "", err
	}

	defer cleanCompose(req, project)
	err = project.Up(ctx, options.Up{})
	if err != nil {
		log.Printf("failed to up: %v", err)
		return false, "", err
	}

	time.Sleep(time.Duration(req.Timeout) * time.Millisecond)

	err = project.Log(ctx, false)
	if err != nil {
		log.Printf("failed to get log: %v", err)
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

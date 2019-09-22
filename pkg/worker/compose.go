package worker

import (
	"context"
	"log"
	"time"

	"github.com/docker/cli/cli/config/configfile"
	clitypes "github.com/docker/cli/cli/config/types"
	"github.com/docker/libcompose/docker"
	"github.com/docker/libcompose/docker/auth"
	dockerctx "github.com/docker/libcompose/docker/ctx"
	"github.com/docker/libcompose/project"
	"github.com/docker/libcompose/project/options"

	pb "gitlab.com/CBCTF/bullseye-runner/proto"
)

func RunDockerCompose(ctx context.Context, req *pb.RunnerRequest) (bool, string, error) {
	yml := req.DockerComposeYml

	configFile := &configfile.ConfigFile{
		AuthConfigs: map[string]clitypes.AuthConfig{
			"localhost:5000": clitypes.AuthConfig{
				// Username: "admin",
				// Password: "password",
				Auth:          req.DockerRegistryToken,
				ServerAddress: "localhost:5000",
			},
		},
	}

	project, err := docker.NewProject(&dockerctx.Context{
		Context: project.Context{
			ComposeBytes: [][]byte{[]byte(yml)},
			ProjectName:  req.Uuid,
		},
		ConfigFile: configFile,
		AuthLookup: auth.NewConfigLookup(configFile),
	}, nil)
	if err != nil {
		return false, "", err
	}

	err = project.Up(ctx, options.Up{})
	if err != nil {
		log.Printf("failed to up: %v", err)
		return false, "", err
	}

	time.Sleep(time.Duration(req.Timeout) * time.Millisecond)

	err = project.Down(ctx, options.Down{})
	if err != nil {
		return false, "", err
	}

	err = project.Kill(ctx, "KILL")
	if err != nil {
		return false, "", err
	}

	return false, "", nil
}

package worker

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
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
	"github.com/lucasjones/reggen"

	pb "gitlab.com/CBCTF/bullseye-runner/proto"
)

const (
	TempDir      = "/tmp"
	FlagSuffix   = "flag"
	SubmitSuffix = "submit"
)

type Runner struct {
	uuid    string
	ctx     context.Context
	req     *pb.RunnerRequest
	project project.APIProject

	flagPath   string
	submitPath string

	x11required  bool
	x11capturing bool
}

func NewRunner(ctx context.Context, req *pb.RunnerRequest) *Runner {
	x11required := req.X11Info != nil
	x11capturing := x11required && req.X11Info.CapExt != ""

	runner := &Runner{
		uuid:         req.Uuid,
		ctx:          ctx,
		req:          req,
		x11required:  x11required,
		x11capturing: x11capturing,
	}

	return runner
}

func (r *Runner) prepareFlags() error {
	flagPath := fmt.Sprintf("%s/%s-%s", TempDir, r.uuid, FlagSuffix)
	submitPath := fmt.Sprintf("%s/%s-%s", TempDir, r.uuid, SubmitSuffix)

	// generate flag from template regex
	flagStr, err := generateFlag(r.req.FlagTemplate)
	if err != nil {
		log.Printf("failed to generate flag: %v", err)
		return err
	}
	log.Printf("generated flag: %s", flagStr)

	// write generated flag
	flagFile, err := os.Create(flagPath)
	if err != nil {
		return err
	}
	_, err = flagFile.WriteString(flagStr)
	if err != nil {
		return err
	}
	if err = flagFile.Close(); err != nil {
		return err
	}

	// create empty flag for submittion
	submitFile, err := os.Create(submitPath)
	if err = submitFile.Close(); err != nil {
		return err
	}

	r.flagPath = flagPath
	r.submitPath = submitPath

	return nil

}

func (r *Runner) cleanFlags() error {
	if err := os.Remove(r.flagPath); err != nil {
		return err
	}
	if err := os.Remove(r.submitPath); err != nil {
		return err
	}
	return nil
}

func (r *Runner) checkFlag() (bool, error) {
	submitBytes, err := ioutil.ReadFile(r.submitPath)
	if err != nil {
		log.Printf("failed to open submitted flag")
		return false, err
	}

	flagBytes, err := ioutil.ReadFile(r.flagPath)
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

func (r *Runner) getNeworkID() string {
	return fmt.Sprintf("%s_default", r.uuid)
}

func (r *Runner) prepareNetwork() error {
	networkID := r.getNeworkID()
	client, err := client.NewEnvClient()
	if err != nil {
		log.Printf("failed to create env client: %v", err)
		return err
	}
	_, err = client.NetworkCreate(r.ctx, networkID, apitypes.NetworkCreate{
		Internal: true,
	})
	if err != nil {
		log.Printf("failed to create network: %v", err)
		return err
	}

	return nil
}

func (r *Runner) cleanNetwork() error {
	networkID := r.getNeworkID()
	client, err := client.NewEnvClient()
	if err != nil {
		return err
	}
	if err := client.NetworkRemove(context.Background(), networkID); err != nil {
		return err
	}

	return nil
}

func (r *Runner) cleanCompose() {
	r.project.Delete(context.Background(), options.Delete{
		RemoveVolume:  true,
		RemoveRunning: true,
	})
	r.cleanNetwork()
}

func (r *Runner) Run() (bool, error) {
	if err := r.prepareFlags(); err != nil {
		return false, err
	}
	defer r.cleanFlags()

	var yml bytes.Buffer
	tpl, err := template.New("yml").Parse(r.req.Yml)
	if err != nil {
		return false, err
	}

	dict := map[string]string{
		"registryHost": r.req.RegistryHost,
		"flagPath":     r.flagPath,
		"submitPath":   r.submitPath,
	}

	err = tpl.Execute(&yml, dict)
	if err != nil {
		return false, err
	}

	log.Printf("execute yml: %#v", yml.String())

	configFile := &configfile.ConfigFile{
		AuthConfigs: map[string]clitypes.AuthConfig{
			r.req.RegistryHost: clitypes.AuthConfig{
				Username: r.req.RegistryUsername,
				Password: r.req.RegistryPassword,
				// Auth:          req.DockerRegistryToken,
				ServerAddress: r.req.RegistryHost,
			},
		},
	}

	project, err := docker.NewProject(&dockerctx.Context{
		Context: project.Context{
			ComposeBytes: [][]byte{yml.Bytes()},
			ProjectName:  r.uuid,
		},
		ConfigFile: configFile,
		AuthLookup: auth.NewConfigLookup(configFile),
	}, nil)
	if err != nil {
		return false, err
	}
	r.project = project

	// create network in advance to make evaluation faster
	if err := r.prepareNetwork(); err != nil {
		return false, err
	}
	defer r.cleanNetwork()

	err = project.Up(r.ctx, options.Up{})
	if err != nil {
		log.Printf("failed to up: %v", err)
		return false, err
	}
	defer r.cleanCompose()

	time.Sleep(time.Duration(r.req.Timeout) * time.Millisecond)

	err = project.Log(r.ctx, false)
	if err != nil {
		log.Printf("failed to get log: %v", err)
		return false, err
	}

	ok, err := r.checkFlag()
	if err != nil {
		return false, err
	}

	if ok {
		return true, nil
	}

	return false, nil
}

func (r *Runner) DryRun() error {
	if err := r.prepareFlags(); err != nil {
		return err
	}
	defer r.cleanFlags()

	var yml bytes.Buffer
	tpl, err := template.New("yml").Parse(r.req.Yml)
	if err != nil {
		return err
	}

	dict := map[string]string{
		"registryHost": r.req.RegistryHost,
		"flagPath":     r.flagPath,
		"submitPath":   r.submitPath,
	}

	err = tpl.Execute(&yml, dict)
	if err != nil {
		return err
	}

	log.Printf("execute yml: %#v", yml.String())

	configFile := &configfile.ConfigFile{
		AuthConfigs: map[string]clitypes.AuthConfig{
			r.req.RegistryHost: clitypes.AuthConfig{
				Username: r.req.RegistryUsername,
				Password: r.req.RegistryPassword,
				// Auth:          req.DockerRegistryToken,
				ServerAddress: r.req.RegistryHost,
			},
		},
	}

	project, err := docker.NewProject(&dockerctx.Context{
		Context: project.Context{
			ComposeBytes: [][]byte{yml.Bytes()},
			ProjectName:  r.uuid,
		},
		ConfigFile: configFile,
		AuthLookup: auth.NewConfigLookup(configFile),
	}, nil)
	if err != nil {
		return err
	}

	err = project.Pull(r.ctx)
	if err != nil {
		return err
	}

	return nil
}

func generateFlag(template string) (string, error) {
	g, err := reggen.NewGenerator(template)
	if err != nil {
		return "", err
	}
	return g.Generate(10), nil
}

func trim(s []byte) string {
	return strings.Trim(fmt.Sprintf("%s", s), " \x00\r\n")
}

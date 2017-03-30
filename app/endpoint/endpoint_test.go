package endpoint

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	docker "github.com/fsouza/go-dockerclient"
	"github.com/ihsw/the-matrix/app/repo"
	"github.com/ihsw/the-matrix/app/resource"
	"github.com/ihsw/the-matrix/app/simpledocker"
)

const defaultResourceName = "db"
const defaultRepoName = "es-bueno"

func cleanResources(t *testing.T, resources resource.Resources) {
	if err := resources.Clean(); err != nil {
		t.Errorf("Could not clean resources: %s", err.Error())
		return
	}
}

func cleanEndpoint(t *testing.T, e Endpoint) {
	if err := e.Clean(); err != nil {
		t.Errorf("Could not clean endpoint: %s", err.Error())
		return
	}
}

func TestNewEndpoint(t *testing.T) {
	dockerClient, err := docker.NewClientFromEnv()
	if err != nil {
		t.Errorf("Could not create a new docker client: %s", err.Error())
		return
	}
	client := simpledocker.NewClient(dockerClient)

	/**
	 * endpoint resources
	 */
	// creating the resource dir
	cwd, err := os.Getwd()
	if err != nil {
		t.Errorf("Could not get working dir: %s", err.Error())
		return
	}
	resourceDir, err := filepath.Abs(fmt.Sprintf("%s/../../%s", cwd, defaultResourceName))
	if err != nil {
		t.Errorf("Could not get absolute filepath for default resource name: %s", err.Error())
		return
	}

	// creating the test resource
	endpointResources, err := resource.NewResources(client, []resource.Opts{resource.Opts{
		Name:                 defaultResourceName,
		DockerfileContextDir: resourceDir,
		EndpointEnvVars:      map[string]string{"DATABASE_HOST": "Db"},
	}})
	if err != nil {
		t.Errorf("Could not create a new resource with default resource %s: %s", defaultResourceName, err.Error())
		return
	}
	defer cleanResources(t, endpointResources)

	/**
	 * endpoint
	 */
	// creating the resource repo
	repo, err := repo.NewRepo(defaultRepoName, client)
	if err != nil {
		t.Errorf("Could not create new repo %s: %s", defaultRepoName, err.Error())
	}

	// creating the endpoint
	endpoint, err := NewEndpoint(repo, endpointResources)
	if err != nil {
		t.Errorf("Could not create a new endpoint based on repo %s: %s", repo.Name, err.Error())
		return
	}
	defer cleanEndpoint(t, endpoint)

	// verifying that it is running
	isRunning, err := client.IsRunning(endpoint.Container)
	if err != nil {
		t.Errorf("Could not check if endpoint container is running: %s", err.Error())
		return
	}
	if !isRunning {
		t.Errorf("Endpoint container %s was not up", endpoint.Name)
		return
	}
}

package resource

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/fsouza/go-dockerclient"
	"github.com/ihsw/the-matrix/app/simpledocker"
)

// NewResources - generates a new list of resources
func NewResources(simpleDocker simpledocker.Client, names map[string]string) ([]Resource, error) {
	resources := []Resource{}
	for name, endpointTarget := range names {
		resource, err := newResource(name, endpointTarget, simpleDocker)
		if err != nil {
			return []Resource{}, err
		}

		resources = append(resources, resource)
	}

	return resources, nil
}

func newResource(name string, endpointTarget string, simpleDocker simpledocker.Client) (Resource, error) {
	resource := Resource{
		Name:           name,
		EndpointTarget: endpointTarget,
		simpleDocker:   simpleDocker,
	}

	var err error
	resource.Container, err = getContainer(resource)
	if err != nil {
		return Resource{}, err
	}

	return resource, nil
}

func getContainer(r Resource) (*docker.Container, error) {
	containerID := fmt.Sprintf("%s-resource", r.Name)
	container, err := r.simpleDocker.GetContainer(containerID)
	if err == nil {
		return container, nil
	}

	log.WithFields(log.Fields{
		"name": r.Name,
	}).Info("Creating resource container")
	container, err = r.simpleDocker.CreateContainer(
		containerID,
		fmt.Sprintf("ihsw/the-matrix-%s", r.Name),
		[]string{},
	)
	if err != nil {
		return nil, err
	}

	if err := r.simpleDocker.StartContainer(container, []string{}); err != nil {
		return nil, err
	}

	return container, nil
}

// Resource - a container for each Endpoint to use (database, etc)
type Resource struct {
	Name           string
	EndpointTarget string
	simpleDocker   simpledocker.Client
	Container      *docker.Container
}

// Clean - stops and removes the Resource's container
func (r Resource) Clean() error {
	if err := r.simpleDocker.StopContainer(r.Container); err != nil {
		return err
	}

	if err := r.simpleDocker.RemoveContainer(r.Container); err != nil {
		return err
	}

	return nil
}

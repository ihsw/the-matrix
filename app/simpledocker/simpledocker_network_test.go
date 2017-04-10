package simpledocker

import (
	"testing"

	"time"

	docker "github.com/fsouza/go-dockerclient"
	"github.com/ihsw/the-matrix/app/util"
)

const defaultTestNetworkName = "test-network"
const defaultNetworkDriver = "bridge"

func createTestNetwork(client Client, namePrefix string, driver string) (*docker.Network, error) {
	name, err := util.GetPrefixedUUID(namePrefix)
	if err != nil {
		return nil, err
	}

	network, err := client.CreateNetwork(name, driver)
	if err != nil {
		return nil, err
	}

	return network, nil
}

func cleanupNetwork(t *testing.T, client Client, network *docker.Network) {
	if err := client.RemoveNetwork(network); err != nil {
		t.Errorf("Could not remove network: %s", err.Error())
		return
	}
}

func TestCreateNetwork(t *testing.T) {
	// creating a simpledocker client
	dockerClient, err := docker.NewClientFromEnv()
	if err != nil {
		t.Errorf("Could not create a new docker client: %s", err.Error())
		return
	}
	client := NewClient(dockerClient)

	network, err := createTestNetwork(client, defaultTestNetworkName, defaultNetworkDriver)
	if err != nil {
		t.Errorf("Could not create network: %s", err.Error())
		return
	}
	defer cleanupNetwork(t, client, network)
}

func TestConnect(t *testing.T) {
	dockerClient, err := docker.NewClientFromEnv()
	if err != nil {
		t.Errorf("Could not create a new docker client: %s", err.Error())
		return
	}
	client := NewClient(dockerClient)

	network, err := createTestNetwork(client, defaultTestNetworkName, defaultNetworkDriver)
	if err != nil {
		t.Errorf("Could not create network: %s", err.Error())
		return
	}
	defer cleanupNetwork(t, client, network)

	hasImage, err := client.HasImage(defaultDbImage)
	if err != nil {
		t.Errorf("Could not check if image exists: %s", err.Error())
		return
	}
	if !hasImage {
		if err := client.PullImage(defaultDbImage, defaultTestImageTag); err != nil {
			t.Errorf("Could not pull image: %s", err.Error())
			return
		}
	}

	_, container, err := createTestContainer(client, defaultTestContainerName, defaultDbImage, nil)
	if err != nil {
		t.Errorf("Could not create container: %s", err.Error())
		return
	}
	defer cleanupContainer(t, client, container)
	if err := client.StartContainer(container, []string{}); err != nil {
		t.Errorf("Could not start container: %s", err.Error())
		return
	}

	time.Sleep(5 * time.Second)
	isRunning, err := client.IsRunning(container)
	if err != nil {
		t.Errorf("Could not check if container was still running: %s", err.Error())
		return
	}
	if !isRunning {
		t.Errorf("Container was not still up after 5s")
		return
	}

	network, err = client.Connect(network, container)
	if err != nil {
		t.Errorf("Could not connect container to network: %s", err.Error())
		return
	}
	if len(network.Containers) != 1 {
		t.Errorf("Could not validate that conatiner was connected to network")
		return
	}
}

func TestCreateContainerWithNetwork(t *testing.T) {
	dockerClient, err := docker.NewClientFromEnv()
	if err != nil {
		t.Errorf("Could not create a new docker client: %s", err.Error())
		return
	}
	client := NewClient(dockerClient)

	network, err := createTestNetwork(client, defaultTestNetworkName, defaultNetworkDriver)
	if err != nil {
		t.Errorf("Could not create network: %s", err.Error())
		return
	}
	defer cleanupNetwork(t, client, network)

	hasImage, err := client.HasImage(defaultDbImage)
	if err != nil {
		t.Errorf("Could not check if image exists: %s", err.Error())
		return
	}
	if !hasImage {
		if err := client.PullImage(defaultDbImage, defaultTestImageTag); err != nil {
			t.Errorf("Could not pull image: %s", err.Error())
			return
		}
	}

	_, container, err := createTestContainer(client, defaultTestContainerName, defaultDbImage, network)
	if err != nil {
		t.Errorf("Could not create container: %s", err.Error())
		return
	}
	defer cleanupContainer(t, client, container)
	if err := client.StartContainer(container, []string{}); err != nil {
		t.Errorf("Could not start container: %s", err.Error())
		return
	}

	time.Sleep(5 * time.Second)
	isRunning, err := client.IsRunning(container)
	if err != nil {
		t.Errorf("Could not check if container was still running: %s", err.Error())
		return
	}
	if !isRunning {
		t.Errorf("Container was not still up after 5s")
		return
	}

	network, err = client.GetNetwork(network.ID)
	if err != nil {
		t.Errorf("Could not fetch network network: %s", err.Error())
		return
	}
	if len(network.Containers) != 1 {
		t.Errorf("Could not validate that conatiner was connected to network")
		return
	}
}

func TestGetContainerIP(t *testing.T) {
	dockerClient, err := docker.NewClientFromEnv()
	if err != nil {
		t.Errorf("Could not create a new docker client: %s", err.Error())
		return
	}
	client := NewClient(dockerClient)

	network, err := createTestNetwork(client, defaultTestNetworkName, defaultNetworkDriver)
	if err != nil {
		t.Errorf("Could not create network: %s", err.Error())
		return
	}
	defer cleanupNetwork(t, client, network)

	hasImage, err := client.HasImage(defaultDbImage)
	if err != nil {
		t.Errorf("Could not check if image exists: %s", err.Error())
		return
	}
	if !hasImage {
		if err := client.PullImage(defaultDbImage, defaultTestImageTag); err != nil {
			t.Errorf("Could not pull image: %s", err.Error())
			return
		}
	}

	_, container, err := createTestContainer(client, defaultTestContainerName, defaultDbImage, network)
	if err != nil {
		t.Errorf("Could not create container: %s", err.Error())
		return
	}
	defer cleanupContainer(t, client, container)
	if err := client.StartContainer(container, []string{}); err != nil {
		t.Errorf("Could not start container: %s", err.Error())
		return
	}

	time.Sleep(5 * time.Second)
	isRunning, err := client.IsRunning(container)
	if err != nil {
		t.Errorf("Could not check if container was still running: %s", err.Error())
		return
	}
	if !isRunning {
		t.Errorf("Container was not still up after 5s")
		return
	}

	ip, err := client.GetContainerIP(network, container)
	if err != nil {
		t.Errorf("Could not get container IP: %s", err.Error())
		return
	}
	if ip == nil {
		t.Errorf("Container IP was nil")
		return
	}
}

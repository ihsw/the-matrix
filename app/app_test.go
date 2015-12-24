package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/fsouza/go-dockerclient"
	"github.com/stretchr/testify/assert"
	"testing"
)

func fail(t *testing.T, err error) {
	assert.Equal(t, err.Error(), nil)
}

func TestTestSuite(t *testing.T) {
	// misc
	repoManager := repoManager{cloneDestination: "./repos"}
	repoNames := []string{
		"omega-jazz",
		"pho-sho",
		"go-home",
		"py-lyfe",
		"es-bueno",
	}
	client, err := docker.NewClient("unix:///var/run/docker.sock")
	if err != nil {
		fail(t, err)
		return
	}

	// spinning up some workers
	reposToDo := make(chan repo)
	reposFinished := runTasks(reposToDo, client)

	// starting it up
	go func() {
		for _, repoName := range repoNames {
			reposToDo <- repo{manager: repoManager, name: repoName, client: client}
		}
		close(reposToDo)
	}()

	// waiting for it to drain out
	for task := range reposFinished {
		if err := task.err; err != nil {
			if err := task.repo.cleanup(); err != nil {
				fail(t, err)
				return
			}
			if err := repoManager.cleanup(); err != nil {
				fail(t, err)
				return
			}
			fail(t, err)
			return
		}
	}

	log.Info("Manager cleanup")
	if err := repoManager.cleanup(); err != nil {
		fail(t, err)
		return
	}
}

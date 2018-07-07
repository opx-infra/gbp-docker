package main

import (
	"fmt"
	"os"

	"github.com/ahmetalpbalkan/dexec"
	docker "github.com/fsouza/go-dockerclient"
	log "github.com/sirupsen/logrus"
)

// BuildConfig contains options for building a Debian package
type BuildConfig struct {
	UID           string
	GID           string
	BuildPath     string
	Distribution  string
	Architecture  string
	Sources       string
	ContainerName string
}

// BuildPackage builds the Debian package using gbp in a container
func (b BuildConfig) BuildPackage() error {
	log.Debug("Getting docker client")
	client, err := docker.NewClientFromEnv()
	if err != nil {
		return err
	}

	log.Debug("Configuring container")
	method, _ := dexec.ByCreatingContainer(docker.CreateContainerOptions{
		Name: b.ContainerName,
		Config: &docker.Config{
			Env: []string{
				"DIST=" + b.Distribution,
				"ARCH=" + b.Architecture,
				"UID=" + b.UID,
				"GID=" + b.GID,
				"EXTRA_SOURCES=" + b.Sources,
			},
			Image: "opxhub/gbp",
		},
		HostConfig: &docker.HostConfig{
			Binds: []string{os.ExpandEnv("$PWD") + ":/mnt"},
		},
	})

	log.Debug("Setting container command")
	execClient := dexec.Docker{Client: client}
	cmd := execClient.Command(method, "buildpackage", b.BuildPath)
	log.Debug("Setting container stdout/stderr to our stdout/stderr")
	cmd.Stderr = os.Stdout
	cmd.Stdout = os.Stdout

	log.WithFields(log.Fields{
		"cmd":  cmd.Path,
		"args": cmd.Args,
	}).Debug("Running")
	if err := cmd.Start(); err != nil {
		return err
	}
	if err := cmd.Wait(); err != nil {
		if exiterr, ok := err.(*dexec.ExitError); ok {
			return fmt.Errorf("process finished with return code %d", exiterr.ExitCode)
		}
		return err
	}
	log.Debug("Container exited successfully")

	return nil
}

// RemoveContainer forcefully removes any running containers for this BuildConfig
func (b BuildConfig) RemoveContainer() error {
	log.Debug("Getting docker client")
	client, err := docker.NewClientFromEnv()
	if err != nil {
		return err
	}

	containers, err := client.ListContainers(docker.ListContainersOptions{
		All:     true,
		Filters: map[string][]string{"name": []string{b.ContainerName}},
	})
	if err != nil {
		return err
	}

	for _, c := range containers {
		log.WithFields(log.Fields{
			"name": c.Names[0],
		}).Info("Removing container")

		err := client.RemoveContainer(docker.RemoveContainerOptions{
			ID:    c.ID,
			Force: true,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

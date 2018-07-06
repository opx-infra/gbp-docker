package main

import (
	"fmt"
	"os"
	"os/signal"
	"os/user"
	"path"

	"github.com/ahmetalpbalkan/dexec"
	docker "github.com/fsouza/go-dockerclient"
	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
)

var buildConfig BuildConfig

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

func main() {
	// SIGINT handler, remove container
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt)
	go func() {
		sig := <-sigc
		switch sig {
		case os.Interrupt:
			log.Debug("Interrupt received.")
			if err := buildConfig.RemoveContainer(); err != nil {
				log.Fatal(err)
			}
		}
	}()

	// Build package and report any errors
	if err := buildConfig.BuildPackage(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	// Set custom container name
	buildConfig.ContainerName = fmt.Sprintf(
		"%s-dbp-%s", os.ExpandEnv("$USER"), path.Base(os.ExpandEnv("$PWD")),
	)

	// Usage message
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `%s

Usage: dbp src/

Builds a Debian package using a Docker container.
Artifacts are found in pool/stretch-amd64/src/

Options:
`, os.Args[0])
		flag.PrintDefaults()
	}

	// Parse CLI arguments / environment
	currentUser, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	flag.StringVarP(&buildConfig.Architecture, "architecture", "a", "amd64", "Debian architecture")
	flag.StringVarP(&buildConfig.Distribution, "distribution", "d", "stretch", "Debian distribution")
	flag.StringVarP(&buildConfig.UID, "uid", "u", currentUser.Uid, "User ID")
	flag.StringVarP(&buildConfig.GID, "gid", "g", currentUser.Gid, "Group ID")
	flag.StringVarP(&buildConfig.Sources, "sources", "s", os.ExpandEnv("$EXTRA_SOURCES"),
		"Extra sources to pull build dependencies from")
	verbose := flag.BoolP("verbose", "v", false, "Print debug messages")

	flag.Parse()

	// Single optional positional argument for which directory to build
	if flag.NArg() == 0 {
		buildConfig.BuildPath = "."
	} else {
		buildConfig.BuildPath = flag.Arg(0)
	}

	if *verbose {
		log.SetLevel(log.DebugLevel)
		log.Debug("Verbosity set to Debug")
	}
}

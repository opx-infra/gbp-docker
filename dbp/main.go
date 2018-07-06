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

var name = fmt.Sprintf(
	"%s-dbp-%s", os.ExpandEnv("$USER"), path.Base(os.ExpandEnv("$PWD")),
)

var (
	arch      string
	dist      string
	uid       string
	gid       string
	sources   string
	buildPath string
)

func main() {
	log.Debug("Getting docker client")
	client, err := docker.NewClientFromEnv()
	if err != nil {
		log.Fatal(err)
	}

	log.Debug("Configuring container")
	method, _ := dexec.ByCreatingContainer(docker.CreateContainerOptions{
		Name: name,
		Config: &docker.Config{
			Env: []string{
				"DIST=" + dist,
				"ARCH=" + arch,
				"UID=" + uid,
				"GID=" + gid,
				"EXTRA_SOURCES=" + sources,
			},
			Image: "opxhub/gbp",
		},
		HostConfig: &docker.HostConfig{
			Binds: []string{os.ExpandEnv("$PWD") + ":/mnt"},
		},
	})

	log.Debug("Setting container command")
	execClient := dexec.Docker{Client: client}
	cmd := execClient.Command(method, "buildpackage", buildPath)
	log.Debug("Setting container stdout/stderr to our stdout/stderr")
	cmd.Stderr = os.Stdout
	cmd.Stdout = os.Stdout

	log.WithFields(log.Fields{
		"cmd":  cmd.Path,
		"args": cmd.Args,
	}).Debug("Running")
	if err := cmd.Start(); err != nil {
		log.Fatal(err.Error())
	}
	if err := cmd.Wait(); err != nil {
		if exiterr, ok := err.(*dexec.ExitError); ok {
			log.WithFields(log.Fields{
				"code": exiterr.ExitCode,
			}).Fatal("Build failed")
		} else {
			log.Fatal(err.Error())
		}
	}
	log.Debug("Container exited successfully")
}

func init() {
	// Parse CLI arguments / environment
	var currentUser, err = user.Current()
	if err != nil {
		log.Fatal(err)
	}
	flag.StringVarP(&arch, "architecture", "a", "amd64", "Debian architecture")
	flag.StringVarP(&dist, "distribution", "d", "stretch", "Debian distribution")
	flag.StringVarP(&uid, "uid", "u", currentUser.Uid, "User ID")
	flag.StringVarP(&gid, "gid", "g", currentUser.Gid, "Group ID")
	flag.StringVarP(&sources, "sources", "s", os.ExpandEnv("$EXTRA_SOURCES"),
		"Extra sources to pull build dependencies from")

	verbose := flag.BoolP("verbose", "v", false, "Print debug messages")

	// Custom usage function
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `%s

Usage: dbp src/

Builds a Debian package using a Docker container.
Artifacts are found in pool/stretch-amd64/src/

Options:
`, os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	if flag.NArg() == 0 {
		buildPath = "."
	} else {
		buildPath = flag.Arg(0)
	}

	if *verbose {
		log.SetLevel(log.DebugLevel)
		log.Debug("Verbosity set to Debug")
	}

	// Handle SIGINT gracefully in container
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt)
	go func() {
		sig := <-sigc
		switch sig {
		case os.Interrupt:
			client, err := docker.NewClientFromEnv()
			if err != nil {
				log.Fatal(err)
			}

			log.Info("Waiting for container to exit.")

			running, err := client.ListContainers(docker.ListContainersOptions{
				All: true,
				Filters: map[string][]string{
					"name": []string{name},
				},
			})
			if err != nil {
				log.Fatal(err)
			}

			for _, c := range running {
				log.WithFields(log.Fields{
					"name": c.Names[0],
				}).Info("Removing container")
				err = client.RemoveContainer(docker.RemoveContainerOptions{
					ID:    c.ID,
					Force: true,
				})
			}
		}
	}()
}

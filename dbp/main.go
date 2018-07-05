package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"os/user"
	"path"

	"github.com/ahmetalpbalkan/dexec"
	docker "github.com/fsouza/go-dockerclient"
	flag "github.com/ogier/pflag"
)

var name = fmt.Sprintf("%s-dbp-%s", os.ExpandEnv("$USER"), path.Base(os.ExpandEnv("$PWD")))

var (
	client    docker.Client
	arch      string
	dist      string
	uid       string
	gid       string
	sources   string
	buildPath string
)

func main() {
	// Instantiate docker client
	client, err := docker.NewClientFromEnv()
	if err != nil {
		log.Fatal(err)
	}

	// Configure container
	method, _ := dexec.ByCreatingContainer(
		docker.CreateContainerOptions{
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
				AutoRemove: true,
				Binds:      []string{os.ExpandEnv("$PWD") + ":/mnt"},
			},
		},
	)

	// Set container command and attach stdout/stderr
	execClient := dexec.Docker{Client: client}
	cmd := execClient.Command(method, "buildpackage", buildPath)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	// Run the command inside the container
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	cmd.Wait()
}

func init() {
	var currentUser, err = user.Current()
	if err != nil {
		log.Fatal(err)
	}
	flag.StringVarP(&arch, "architecture", "a", "amd64", "Debian architecture")
	flag.StringVarP(&dist, "distribution", "d", "stretch", "Debian distribution")
	flag.StringVarP(&uid, "uid", "u", currentUser.Uid, "User ID")
	flag.StringVarP(&gid, "gid", "g", currentUser.Gid, "Group ID")
	flag.StringVarP(
		&sources,
		"sources",
		"s",
		os.ExpandEnv("$EXTRA_SOURCES"),
		"Extra sources to pull build dependencies from",
	)

	flag.Parse()

	if flag.NArg() == 0 {
		buildPath = "."
	} else {
		buildPath = flag.Arg(0)
	}

	// Handle SIGINT gracefully in container
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt)
	go func() {
		sig := <-sigc
		switch sig {
		case os.Interrupt:
			log.Println(sig, "Waiting for container to exit.")

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
				log.Println("Removing your container", c.Names[0])
				err = client.RemoveContainer(docker.RemoveContainerOptions{
					ID:    c.ID,
					Force: true,
				})
			}
		}
	}()
}

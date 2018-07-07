package main

import (
	"fmt"
	"os"
	"os/signal"
	"os/user"
	"path"

	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
)

func main() {
	buildConfig := newBuildConfigFromArgs()

	// Set custom container name
	buildConfig.ContainerName = fmt.Sprintf(
		"%s-dbp-%s", os.ExpandEnv("$USER"), path.Base(os.ExpandEnv("$PWD")),
	)

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

func newBuildConfigFromArgs() BuildConfig {
	// Our build configuration
	var buildConfig BuildConfig

	// Parse args directly into it
	currentUser, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	flag.StringVarP(&buildConfig.UID, "uid", "u", currentUser.Uid, "User ID")
	flag.StringVarP(&buildConfig.GID, "gid", "g", currentUser.Gid, "Group ID")
	flag.StringVarP(&buildConfig.Architecture, "architecture", "a", "amd64", "Debian architecture")
	flag.StringVarP(&buildConfig.Distribution, "distribution", "d", "stretch", "Debian distribution")
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

	return buildConfig
}

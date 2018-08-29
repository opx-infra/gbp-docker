# gbp-docker

*git-buildpackage + docker*

gbp-docker is an opinionated Debian development environment in a container. A small Python (>= 3.5) script is included to help you manage the container lifecycle. Example `docker run` invocations are also included later in this page.

# dbp

dbp is used to manage the persistence of the development environment, enabling compiler and dependency caches for faster builds. dbp requires Python (>= 3.5).

## Installation

Fetch the latest release.

```bash
curl -LO https://raw.githubusercontent.com/opx-infra/gbp-docker/master/dbp
chmod +x ./dbp
```

## Usage

* `dbp build src/` runs an out-of-tree build and stores build artifacts in `./pool/` for easy publishing
* `dbp shell` launches an interactive bash shell in the development environment container
* `dbp run` starts a persistent container in the background
* `dbp rm` removes the persistent container

Both `dbp build` and `dbp run` use temporary containers if no container exists.

## Advanced usage

Here are some fun things you can do with `dbp`.

### Build against extra apt sources

`dbp` will read from the following list of inputs for extra apt sources. These sources must be in standard sources.list format.

1. `--extra-sources` argument
1. `EXTRA_SOURCES` environment variable
1. `./.extra_sources` file
1. `~/.extra_sources` file

For example, fill `~/.extra_sources` with
```bash
deb     http://deb.openswitch.net/stretch stable opx opx-non-free
deb-src http://deb.openswitch.net/stretch stable opx
```
and `dbp build` will search OpenSwitch for build dependencies.

```bash
$ dbp -v build src
INFO:dbp:Loaded extra sources:
deb     http://deb.openswitch.net/stretch stable opx opx-non-free
deb-src http://deb.openswitch.net/stretch stable opx
```

### Build a single package in a non-persistent container

```bash
dbp build src
```

* Builds artifacts for the default Debian distribution
* Uses packages found in `./pool/stretch-amd64` as build dependencies
* Deposits artifacts in `./pool/stretch-amd64/src/`
* If workspace container does not exist, a container is created for this build and destroyed after
* If the workspace container already exists, it is used for the build and *not* destroyed after

```bash
dbp --dist buster build src
```

* Builds the package against Buster
* Deposits artifacts in `./pool/buster-amd64/src/`

```bash
EXTRA_SOURCES="deb http://deb.openswitch.net/stretch 3.0.0 opx opx-non-free"
dbp build src
```

* Adds `EXTRA_SOURCES` to `sources.list`

### Develop inside a persistent development container

Using the `run` subcommand launches a persistent development container. This container will only be explicitly removed when `dbp rm` is run in the same directory. Then use `dbp shell` to enter this container.

```bash
dbp run
dbp shell

# Now we are inside the container (denoted by $ prompt)

# Install build dependencies and build the package
$ cd src/
$ gbp buildpackage

# On failed builds, avoid the long gbp build time by quickly rebuilding
$ fakeroot debian/rules build

# Manually clean up
$ fakeroot debian/rules clean

# Add a new source for build dependencies by appending to the env var
$ export EXTRA_SOURCES="$EXTRA_SOURCES
deb http://deb.openswitch.net/stretch 3.0.0 opx opx-non-free"

# Run gbp buildpackage again to do a clean build and install any new dependencies
$ gbp buildpackage

# Exit the container
$ exit

# Remove the container when finished (or use `dbp run` again to re-enter the same container)
dbp rm
```

### Pass additional `git-buildpackage` options

For example, skip building when tagging by passing the correct flag.

```bash
dbp build src --gbp="--git-tag-only"
```

### Pull any Docker image updates

```bash
dbp pull
dbp -d buster pull
```

# Docker run commands

### Build a single package in a container

```bash
docker run --rm -it -v "$(pwd):/mnt" -e UID=$(id -u) -e GID=$(id -g) -e EXTRA_SOURCES opxhub/gbp:stretch build ./src/
```

### Develop in a development container

```bash
docker run --rm -it -v "$(pwd):/mnt" -e UID=$(id -u) -e GID=$(id -g) -e EXTRA_SOURCES opxhub/gbp:stretch-dev
```

# Building the container images

First, generate the Dockerfiles from the templates.

```bash
make update
```

Next, build the images.

```bash
make DIST=stretch
make DIST=buster
```

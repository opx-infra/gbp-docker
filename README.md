# gbp-docker

*git-buildpackage + docker*

[![Build status](https://badge.buildkite.com/5dbfd1f5cf0ff9311fd6543a5ac976de409fbc8cdd6ecee299.svg)](https://buildkite.com/opx/opx-infra-gbp-docker)

This is the Git repository of the "official" OpenSwitch build and development environments.

## Quick start

```bash
docker run -v "$(pwd):/mnt" -e UID=$(id -u) -e GID=$(id -g) opxhub/gbp:stretch build ./src/
```

Build artifacts can be found in `./pool/stretch-amd64/src/`.

## Building software

The build variant of the image builds the package in a separate directory and exports the results to `./pool/stretch-amd64/src/`.

This alias will make building Debian packages a breeze.

```bash
alias dbp='docker run --rm -it -v "$(pwd):/mnt" -e UID=$(id -u) -e GID=$(id -g) -e EXTRA_SOURCES opxhub/gbp:stretch build'
```

Use it like this.

```bash
dbp ./src/
```

## Developing software

Use the development image variant for interactive sessions. It contains helpful tools, such as `vim`. When building in this image, **the build is done directly in the source tree and artifacts are deposited in the parent directory**.

```bash
# Enter a development container
docker run --rm -it -v "$(pwd):/mnt" -e UID=$(id -u) -e GID=$(id -g) -e EXTRA_SOURCES opxhub/gbp:stretch-dev

# Now we are inside the container (denoted by $ prompt)

# Install build dependencies and build the package
$ cd src/
$ gbp buildpackage

# On failed builds, avoid the long gbp build time by quickly rebuilding
$ fakeroot debian/rules build

# Create a test binary after fixing a failed build
$ fakeroot debian/rules binary

# Manually clean up
$ fakeroot debian/rules clean

# Run gbp buildpackage again to do a clean build and install any new dependencies
$ gbp buildpackage
```

## Build options

### Building against different Debian distributions

Build against different Debian distributions by changing the image tag.

Examples:

- Stretch: Use `opxhub/gbp:stretch` and `opxhub/gbp:stretch-dev`

### Building against custom Apt sources

Set the `EXTRA_SOURCES` environment variable like any `sources.list` file. This work for both the build and develop variants

```bash
export EXTRA_SOURCES="deb http://deb.openswitch.net/stretch stable main opx opx-non-free"
dbp ./src/
```

# Building the images

First, generate the Dockerfiles from the templates.

```bash
make update
```

Next, build the images.

```bash
make stretch
```

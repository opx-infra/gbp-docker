# gbp-docker

*git-buildpackage + docker*

gbp-docker is an opinionated Debian build and development environment in a container.

## Quick start (build)

Given a directory named `./src` with a `debian/` directory, build it with

```bash
docker run -it --rm -v "$(pwd):/mnt" -e UID=$(id -u) -e GID=$(id -g) opxhub/gbp:stretch bash -c 'cd src/; gbp buildpackage'
```

## Quick start (develop)

Launch a development container with

```bash
docker run -it --rm \
  --name=${USER}-dbp-$(basename $(pwd)) \
  --hostname=stretch \
  -v "$(pwd):/mnt" \
  -v "$HOME/.gitconfig:/etc/skel/.gitconfig:ro" \
  -e UID=$(id -u) \
  -e GID=$(id -g) \
  -e EXTRA_SOURCES \
  opxhub/gbp:stretch-dev
```

## Build against extra apt sources

git-buildpackage will use load `$EXTRA_SOURCES` for build dependencies if specified. Use the same format as `/etc/apt/sources.list`.

```bash
export EXTRA_SOURCES="
deb     http://deb.openswitch.net/stretch stable opx opx-non-free
deb-src http://deb.openswitch.net/stretch stable opx"
```

## Pool packages for publishing

A script is provided which will pool packages into respective directories. Simply run it on every changes file.

```bash
for f in *.changes; do pool-packages $f; done
```

Packages will be found in `./pool/${DIST}-${ARCH}/src`.

# Building the container images

Generate the Dockerfiles from the templates and build the images.

```bash
make DIST=stretch
make DIST=buster
make DIST=xenial
make DIST=bionic
```

# gbp-docker

*git-buildpackage + docker*

gbp-docker is an opinionated Debian build and development environment in a container.

## Quick start (build)

Given a directory named `./src` with a `debian/` directory, build it with

```bash
docker run -it --rm -v "$(pwd):/mnt" -e UID=$(id -u) -e GID=$(id -g) opxhub/gbp:stretch build src
```

Build artifacts are deposited in `./pool/stretch-amd64/src/` for easy sharing.

## Quick start (develop)

Given a directory named `./src` with a `debian/` directory, launch a development container with

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

Inside, make use of `git-buildpackage` and `debian/rules` to quickly rebuild while developing

```bash
cd src/
gbp buildpackage

# On failed builds, avoid the long gbp build time by quickly rebuilding
fakeroot debian/rules build

# Manually clean up
fakeroot debian/rules clean

# Run an "official" build (what the CI runs)
build
```

To use build artifacts for future builds, either build with `build` or manually copy the artifacts into `./pool/stretch-amd64/src/`.

## Advanced usage

Here are some fun things you can do with gbp-docker.

### Build against extra apt sources

git-buildpackage will use load `$EXTRA_SOURCES` for build dependencies if specified.

For example, build against OPX packages with this variable

```bash
export DIST=stretch
export EXTRA_SOURCES="
deb     http://deb.openswitch.net/stretch stable opx opx-non-free
deb-src http://deb.openswitch.net/stretch stable opx
"
```

and git-buildpackage will search OpenSwitch for build dependencies.

```bash
docker run -it --rm \
  --name=${USER}-dbp-$(basename $(pwd)) \
  --hostname=${DIST} \
  -v "$(pwd):/mnt" \
  -v "$HOME/.gitconfig:/etc/skel/.gitconfig:ro" \
  -e UID=$(id -u) \
  -e GID=$(id -g) \
  -e EXTRA_SOURCES \
  opxhub/gbp:${DIST}-dev
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

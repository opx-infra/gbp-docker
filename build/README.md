# Docker Builder for Debian Packages

## Quick Start

```bash
docker run --rm -it \
  -v "$(pwd):/mnt" \
  -e DIST=stretch \
  -e ARCH=amd64 \
  -e UID=$(id -u) \
  -e GID=$(id -g) \
  opxhub/gbp buildpackage src/
```

Build artifacts are found in `pool/${DIST}-${ARCH}/src/`.

## Recommended alias

This alias may help with the long command. Environment variables should be set beforehand.

```bash
alias dbp='docker run --rm -it -v "$(pwd):/mnt" -e DIST -e ARCH -e EXTRA_SOURCES -e UID=$(id -u) -e GID=$(id -g) opxhub/gbp buildpackage'
```

Usage:

```bash
DIST=stretch ARCH=amd64 dbp src/
```

## Adding additional package repositories

This example uses the alias from the previous step.

```bash
$ export EXTRA_SOURCES="
deb     http://deb.openswitch.net/stretch unstable opx opx-non-free
deb-src http://deb.openswitch.net/stretch unstable opx
"
$ DIST=stretch ARCH=amd64 dbp src/
```

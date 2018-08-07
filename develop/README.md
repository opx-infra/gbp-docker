# Docker Builder for Debian Packages

## Quick Start

```bash
docker run --rm -it \
  -v "$(pwd):/mnt" \
  -v "$HOME/.gitconfig:/home/opx/.gitconfig" \
  -v "/etc/localtime:/etc/localtime:ro" \
  -e DIST=stretch \
  -e ARCH=amd64 \
  -e UID=$(id -u) \
  -e GID=$(id -g) \
  opxhub/gbp-dev
```

Beware root access to ```/mnt``` within the docker.

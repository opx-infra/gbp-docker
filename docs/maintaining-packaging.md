# Maintaining Debian Packaging

This document outlines the creation of a Debian packaging fork. See the [gbp manual](http://honk.sigxcpu.org/projects/git-buildpackage/manual-html/gbp.import.upstream-git.html) for more information.

## Fork creation

For this example, we'll be creating a fork of `pam_tacplus` since it currently fails to build for Debian Stretch.

1. Clone the upstream repository per the gbp manual.

```bash
git clone --no-checkout -o upstream https://github.com/jeroennijhof/pam_tacplus
```

2. Create a release branch for the target Debian distribution from the tag you wish to release.

```bash
git -C pam_tacplus checkout -b debian/stretch v1.4.1
```

3. In our specific example, we need to add a build dependency and change an install location.

```bash
sed -i 's/libpam-dev/libpam-dev, libssl-dev/' pam_tacplus/debian/control
sed -i 's/sbin/bin/' pam_tacplus/debian/libtac2-bin.install
git -C pam_tacplus commit -asm 'fix: Build with git-buildpackage on Stretch'
```

4. Build the package.

```bash
docker run -it --rm -v "$(pwd):/mnt" -e UID=$(id -u) -e GID=$(id -g) opxhub/gbp:stretch \
  build pam_tacplus --git-debian-branch=debian/stretch --git-tag
```

5. Push our branch and tag to the new fork.

```bash
git -C pam_tacplus remote add origin https://github.com/opx-infra/pam_tacplus
git -C pam_tacplus push origin debian/stretch
git -C pam_tacplus push origin --tags
```

## Fork update (new upstream version)

TODO...

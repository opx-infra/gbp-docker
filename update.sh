#!/usr/bin/env bash
set -euo pipefail
#set -x

debians=( stretch buster )
for dist in "${debians[@]}"; do
  mkdir -p "debian/$dist/base"
  sed -r \
      -e 's!%%DEBIAN-DIST%%!'"$dist"'!g' \
      -e 's!%%DEBIAN-ARCH%%!'"amd64"'!g' \
      "Dockerfile-debian-base.template" >"debian/$dist/base/Dockerfile"
  sed -r \
      -e 's!%%DEBIAN-DIST%%!'"$dist"'!g' \
      -e 's!%%DEBIAN-ARCH%%!'"amd64"'!g' \
      "Dockerfile-debian-dev.template" >"debian/$dist/Dockerfile"
done

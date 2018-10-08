#!/usr/bin/env bash
set -euo pipefail
#set -x

debians=( stretch buster )
for dist in "${debians[@]}"; do
  mkdir -p "debian/$dist/base"
  sed -r \
      -e 's!%%OS%%!debian!g' \
      -e 's!%%DIST%%!'"$dist"'!g' \
      -e 's!%%ARCH%%!amd64!g' \
      "Dockerfile-base.template" >"debian/$dist/base/Dockerfile"
  sed -r \
      -e 's!%%OS%%!debian!g' \
      -e 's!%%DIST%%!'"$dist"'!g' \
      -e 's!%%ARCH%%!amd64!g' \
      "Dockerfile-dev.template" >"debian/$dist/Dockerfile"
done

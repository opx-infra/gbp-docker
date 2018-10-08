#!/usr/bin/env bash
set -euo pipefail
#set -x

update() {
  os="$1"
  dist="$2"

  mkdir -p "$os/$dist/base"
  sed -r \
      -e 's!%%OS%%!'"$os"'!g' \
      -e 's!%%DIST%%!'"$dist"'!g' \
      -e 's!%%ARCH%%!amd64!g' \
      "Dockerfile-base.template" >"$os/$dist/base/Dockerfile"
  sed -r \
      -e 's!%%OS%%!'"$os"'!g' \
      -e 's!%%DIST%%!'"$dist"'!g' \
      -e 's!%%ARCH%%!amd64!g' \
      "Dockerfile-dev.template" >"$os/$dist/Dockerfile"
}

os=debian
debians=( stretch buster )
for dist in "${debians[@]}"; do
  update "$os" "$dist"
done

os=ubuntu
ubuntus=( xenial bionic )
for dist in "${ubuntus[@]}"; do
  update "$os" "$dist"
done

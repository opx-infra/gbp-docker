#!/usr/bin/env bash
set -euo pipefail
#set -x

docker push "opxhub/gbp:${DIST}"
docker push "opxhub/gbp:${DIST}-dev"

if [[ -z "$BUILDKITE_TAG" ]]; then
  exit 0
fi

shopt -s extglob
case "$BUILDKITE_TAG" in
  v*.*.*) echo "+++ Publishing release $BUILDKITE_TAG";;
  *) echo "--- Skipping tag $BUILDKITE_TAG"; exit 0;;
esac

docker tag "opxhub/gbp:${DIST}" "opxhub/gbp:${BUILDKITE_TAG}-${DIST}"
docker tag "opxhub/gbp:${DIST}-dev" "opxhub/gbp:${BUILDKITE_TAG}-${DIST}-dev"
docker push "opxhub/gbp:${BUILDKITE_TAG}-${DIST}"
docker push "opxhub/gbp:${BUILDKITE_TAG}-${DIST}-dev"

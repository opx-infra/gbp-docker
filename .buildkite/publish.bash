#!/usr/bin/env bash
set -euo pipefail
#set -x

docker tag "opxhub/gbp:$DIST" "opxhub/gbp:$(git rev-parse HEAD)-${DIST}"
docker tag "opxhub/gbp:${DIST}-dev" "opxhub/gbp:$(git rev-parse HEAD)-${DIST}-dev"

docker push "opxhub/gbp:${DIST}"
docker push "opxhub/gbp:$(git rev-parse HEAD)-${DIST}"
docker push "opxhub/gbp:${DIST}-dev"
docker push "opxhub/gbp:$(git rev-parse HEAD)-${DIST}-dev"

#!/usr/bin/env bash

if [[ -n "${UID}" ]]; then
    GID=${GID:-${UID}}

    if ! grep -q build /etc/group; then
        groupadd --non-unique --gid="${GID}" build
        useradd --non-unique --uid="${UID}" --gid="${GID}" --create-home --shell /bin/bash build
    fi

    exec gosu build "$@"
fi

exec "$@"


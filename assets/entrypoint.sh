#!/usr/bin/env bash

if [[ -n "${UID}" ]]; then
    GID=${GID:-${UID}}

    if ! grep -q build /etc/group; then
        groupadd -o --gid="${GID}" build
        useradd -o --uid="${UID}" --gid="${GID}" -s /bin/bash build
    fi

    exec gosu build "$@"
fi

exec "$@"


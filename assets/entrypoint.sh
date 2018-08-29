#!/usr/bin/env bash

if [[ -n "${UID}" ]]; then
    GID=${GID:-${UID}}

    if ! grep -q build /etc/group; then
        # first launch
        groupadd --non-unique --gid="${GID}" build
        useradd --non-unique --uid="${UID}" --gid="${GID}" --create-home --shell /bin/bash build

        if [[ -n "${EXTRA_SOURCES}" ]]; then
            echo "${EXTRA_SOURCES}" >/etc/apt/sources.list.d/20extra.list
        fi
    fi

    exec gosu build "$@"
fi

exec "$@"


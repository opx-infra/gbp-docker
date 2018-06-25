FROM debian:stretch-slim

RUN apt-get -qq update && apt-get -qq upgrade -y && apt-get -qq install -y \
    git-buildpackage \
    gosu \
 && apt-get -qq autoremove -y && apt-get clean && rm -rf /var/lib/apt/lists/*

# Set up for the user we will create at runtime
RUN mkdir -p /home/build && chmod -R 777 /home/build \
 && echo 'build ALL=(ALL) NOPASSWD:ALL' >>/etc/sudoers \
 && echo '%build ALL=(ALL) NOPASSWD:ALL' >>/etc/sudoers \
 && echo 'Defaults env_keep += "DIST ARCH"' >>/etc/sudoers

COPY assets/apt-preferences /etc/apt/preferences
COPY assets/buildpackage /usr/bin/buildpackage
COPY assets/entrypoint.sh /
COPY assets/install-build-deps /usr/bin/install-build-deps

# Add OpenSwitch GPG key
RUN gpg --keyserver hkp://keyserver.ubuntu.com:80 --recv-key AD5073F1 \
 && gpg --export AD5073F1 >/usr/share/keyrings/opx-archive-keyring.gpg \
 && apt-key add /usr/share/keyrings/opx-archive-keyring.gpg \
 && mkdir -p /etc/apt/sources.list.d/

WORKDIR /mnt
VOLUME /mnt

ENTRYPOINT ["/entrypoint.sh"]
CMD ["buildpackage"]

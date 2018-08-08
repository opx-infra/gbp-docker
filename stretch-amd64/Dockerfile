FROM debian:stretch-slim

ENV ARCH amd64
ENV DIST stretch

RUN apt-get -qq update && apt-get -qq upgrade -y && apt-get -qq install -y \
    git-buildpackage \
    gosu \
    python-pip \
 && apt-get -qq autoremove -y && apt-get clean && rm -rf /var/lib/apt/lists/*

# Pyang not available as Debian package
RUN pip2 install pyang requests-file \
 && ln -s /usr/local/bin/pyang /usr/bin

# Add OpenSwitch GPG key
RUN gpg --keyserver hkp://keyserver.ubuntu.com:80 --recv-key AD5073F1 \
 && gpg --export AD5073F1 >/usr/share/keyrings/opx-archive-keyring.gpg \
 && apt-key add /usr/share/keyrings/opx-archive-keyring.gpg \
 && mkdir -p /etc/apt/sources.list.d/

# Set up for the user we will create at runtime
RUN mkdir -p /home/build && chmod -R 777 /home/build \
 && echo 'build ALL=(ALL) NOPASSWD:ALL' >>/etc/sudoers \
 && echo '%build ALL=(ALL) NOPASSWD:ALL' >>/etc/sudoers \
 && echo 'Defaults env_keep += "DIST ARCH"' >>/etc/sudoers

COPY apt-preferences /etc/apt/preferences
COPY buildpackage /usr/bin/buildpackage
COPY gbp.conf /home/build/.gbp.conf
COPY entrypoint.sh /
COPY install-build-deps /usr/bin/install-build-deps

WORKDIR /mnt
VOLUME /mnt

ENTRYPOINT ["/entrypoint.sh"]
CMD ["buildpackage"]

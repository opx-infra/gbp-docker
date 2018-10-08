FROM ubuntu:xenial

ENV ARCH       amd64
ENV CCACHE_DIR /home/build/ccache
ENV DIST       xenial
ENV LANG       C.UTF-8
ENV LC_ALL     C.UTF-8
ENV PATH       "/usr/lib/ccache:$PATH"
ENV TERM       xterm-256color

# Install packages for building Debian packages
RUN apt-get -qq update && apt-get -qq upgrade -y \
 && DEBIAN_FRONTEND=noninteractive apt-get -qq install -y \
    ccache \
    debhelper \
    equivs \
    git-buildpackage \
    gosu \
    python-pip \
    quilt \
 && apt-get -qq autoremove -y && apt-get clean && rm -rf /var/lib/apt/lists/*

# Pyang not available as Debian package
RUN pip2 install pyang requests-file \
 && ln -s /usr/local/bin/pyang /usr/bin

# Add OpenSwitch GPG key
RUN apt-key adv --fetch-keys http://deb.openswitch.net/opx.asc \
 && mkdir -p /etc/apt/sources.list.d/

# Set up for the user we will create at runtime
RUN echo 'build ALL=(ALL) NOPASSWD:ALL' >>/etc/sudoers \
 && echo '%build ALL=(ALL) NOPASSWD:ALL' >>/etc/sudoers

COPY assets/apt-preferences    /etc/apt/preferences
COPY assets/bash_aliases       /etc/skel/.bash_aliases
COPY assets/build              /usr/bin/build
COPY assets/entrypoint.sh      /entrypoint.sh
COPY assets/gbp.conf           /etc/skel/.gbp.conf
COPY assets/install-build-deps /usr/bin/install-build-deps

WORKDIR /mnt
VOLUME /mnt

ENTRYPOINT ["/entrypoint.sh"]
CMD ["build"]

# vi:syntax=dockerfile

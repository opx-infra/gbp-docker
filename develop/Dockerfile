FROM opxhub/gbp

# add extra tools required for development
RUN apt-get -qq update && apt-get -qq upgrade -y && \
    apt-get -qq install -y vim

# start a shell for the work
CMD ["/bin/bash"]

# Dockerfile for rabbit-herder
#
# To build:
# $ docker build -t micahhausler/rabbit-herder-builder -f Dockerfile.build .
# $ docker run --rm -v $(pwd):/usr/src/app micahhausler/rabbit-herder-builder
# $ docker build -t micahhausler/rabbit-herder .
#
# To run:
# $ docker run micahhausler/rabbit-herder

FROM busybox

MAINTAINER Micah Hausler, <micah.hausler@ambition.com>

COPY rabbit-herder /bin/rabbit-herder
RUN chmod 755 /bin/rabbit-herder

ENTRYPOINT ["/bin/rabbit-herder"]

FROM busybox:latest

ADD ./pmss /usr/bin/pmss

ENTRYPOINT [ "pmss", "server" ]
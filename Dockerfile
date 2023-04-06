FROM scratch

ADD ./pmss /usr/bin/pmss

ENTRYPOINT [ "pmss", "server" ]
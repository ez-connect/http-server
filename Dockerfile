FROM docker.io/library/busybox:latest

COPY dist/http-server-linux .
RUN chmod +x http-server-linux \
	&& mv http-server-linux /bin/http-server

ENTRYPOINT http-server .

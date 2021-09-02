FROM docker.io/library/busybox:latest
MAINTAINER Vinh thanh.vinh@hotmail.com

WORKDIR /app
COPY build/http-server-linux-amd64 .
RUN chmod +x http-server-linux-amd64 && mkdir /public

ENTRYPOINT /app/http-server-linux-amd64
CMD -root /public

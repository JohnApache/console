FROM golang:1.9-alpine AS builder

ENV GOBIN=$GOPATH/bin

RUN go get -u -v github.com/wzshiming/console/cmd/web_console

FROM alpine:3.6

COPY --from=builder $(go env GOBIN)/web_console /usr/local/bin/

VOLUME /var/run/docker.sock
EXPOSE 8888

CMD /usr/local/bin/web_console

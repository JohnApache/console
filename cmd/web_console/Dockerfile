FROM golang:1.9-alpine AS builder

ENV GOBIN=$GOPATH/bin
RUN apk add -U git gcc musl-dev
RUN go get -u -v github.com/wzshiming/console/cmd/web_console

FROM alpine:3.6
WORKDIR /usr/local/bin/

COPY --from=builder $(go env GOBIN)/web_console /usr/local/bin/
COPY --from=builder $(go env GOPATH)/github.com/wzshiming/console/cmd/web_console/static/ /usr/local/bin/

VOLUME /var/run/docker.sock
EXPOSE 8888

CMD /usr/local/bin/web_console

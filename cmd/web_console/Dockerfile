
FROM golang:1.9.3-alpine3.7 AS builder
ENV GOPATH=/tmp/gopath/
ENV GOBIN=$GOPATH/bin/
# RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
RUN apk add -U git gcc musl-dev
RUN go get -u -v github.com/wzshiming/console/cmd/web_console

FROM scratch
COPY --from=builder /tmp/gopath/bin/web_console /
# COPY --from=builder /tmp/gopath/src/github.com/wzshiming/console/cmd/web_console/static/ /static/
WORKDIR /
VOLUME /var/run/docker.sock
EXPOSE 8888
CMD /web_console

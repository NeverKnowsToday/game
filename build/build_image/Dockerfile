FROM getourneau/alpine-golang-bash:latest as build

MAINTAINER <shilixingchen>

ENV GO111MODULE="off"

ENV BUILDPATH=github.com/gopath/game

RUN mkdir -p /go/src/${BUILDPATH}

ADD server.tar.gz /go/src/${BUILDPATH}

WORKDIR /go/src/${BUILDPATH}/server

RUN go build -o game

# 构建运行镜像
FROM alpine:latest

ENV WORKPATH=/opt/gopath/src/github.com/game/server
ENV BUILDPATH=/go/src/github.com/gopath/game/server

RUN mkdir -p ${WORKPATH}

WORKDIR ${WORKPATH}

COPY --from=build ${BUILDPATH}/game ./

COPY --from=build $BUILDPATH/config.yaml ./

COPY --from=build $BUILDPATH/ErrorCue.json ./


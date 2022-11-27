FROM golang:1.14 as build

ENV GOPROXY https://goproxy.io
ENV GO111MODULE on

ENV BUILDPATH=github.com/game

RUN mkdir -p /go/src/${BUILDPATH}

WORKDIR /go/src/${BUILDPATH}

ADD game.tar.gz /go/src/${BUILDPATH}

WORKDIR /go/src/${BUILDPATH}/game/build/

CMD ["/bin/bash", "-c", "./build.sh"]

# 构建运行镜像
FROM golang:1.14

RUN mkdir -p /opt/game

WORKDIR /opt/game

COPY --from=build /go/src/github.com/game/server/game ./
#
#COPY --from=build /go/src/github.com/game/server/config.yaml ./
#
#COPY --from=build /go/src/github.com/game/server/ErrorCue.json ./
#
#RUN mkdir sdk
#
#COPY --from=build /go/src/github.com/peersafe/fabric_gateway/server/sdk/client_sdk.yaml ./sdk/
#
#COPY --from=build /go/src/github.com/peersafe/fabric_gateway/server/sdk/crypto-config ./sdk/


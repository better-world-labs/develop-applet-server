FROM golang:1.19.1-alpine3.16 as builder
WORKDIR /app

RUN sed 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' -i /etc/apk/repositories && \
    apk add make git

ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct

COPY ["go.mod", "go.sum", "Makefile", "./"]

RUN --mount=type=cache,target=/go,id=go go mod download

COPY . .

#编译
RUN --mount=type=cache,target=/go,id=go make build

FROM alpine:3.16
RUN sed 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' -i /etc/apk/repositories && \
    apk update && \
    apk add tzdata && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone

WORKDIR /app
COPY config etc
COPY --from=builder /app/bin/server /app/server

ARG ENVIRONMENT=dev
ENV ENV=${ENVIRONMENT} CONF=/app/etc

CMD ./server
EXPOSE 8080
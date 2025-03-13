# syntax=docker/dockerfile:1
FROM --platform=$BUILDPLATFORM tonistiigi/xx:1.2.1 AS xx
FROM --platform=$BUILDPLATFORM golang:alpine AS build
RUN mkdir -p /var/cache/apk&& \
    ln -s /var/cache/apk /etc/apk/cache
RUN --mount=type=cache,target=/var/cache/apk \
    apk add --no-cache clang lld make automake autoconf bash git pkgconf
COPY --from=xx / /
ARG TARGETPLATFORM
ARG BUILDPLATFORM
ARG XX_TRIPLE
ARG XX_ALPINE_ARCH
RUN --mount=type=cache,target=/var/cache/apk \
    xx-apk add --no-cache musl-dev gcc pkgconf
RUN --mount=type=cache,target=/var/cache/apk \
    xx-apk add --no-cache fuse3 fuse3-dev fuse3-static linux-headers
ENV CGO_ENABLED=1
WORKDIR /app
COPY . /app
ARG CONFIG_ARGS
RUN --mount=type=cache,target=/go/pkg/mod \
    GO=xx-go \
    GOARCH=`xx-info arch` \
    CC=xx-clang \
    TARGET=`xx-clang --print-target-triple` \
    make all $CONFIG_ARGS
FROM scratch
COPY --from=build /app/ll-killer /ll-killer

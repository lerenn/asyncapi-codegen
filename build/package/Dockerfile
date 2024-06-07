# Default arguments
ARG BUILDPLATFORM=linux/amd64

ARG TARGETOS=linux
ARG TARGETARCH=amd64

ARG BUILDBASEIMAGE=golang:alpine
ARG TARGETBASEIMAGE=alpine:latest

# Building image
FROM --platform=${BUILDPLATFORM} ${BUILDBASEIMAGE} AS build

# Disable CGO
ENV CGO_ENABLED 0

# Get all remaining code
RUN mkdir -p /go/src/github.com/TheSadlig/asyncapi-codegen
COPY ./ /go/src/github.com/TheSadlig/asyncapi-codegen

# Set the workdir
WORKDIR /go/src/github.com/TheSadlig/asyncapi-codegen

# Build everything in cmd/
RUN --mount=type=cache,target=/root/.cache/go-build \
    GOOS=${TARGETOS} GOARCH=${TARGETARCH} CGO_ENABLED=0 \ 
    go install ./cmd/*

# Get final base image
FROM --platform=${TARGETOS}/${TARGETARCH} ${TARGETBASEIMAGE}

# Get binaries
COPY --from=build /go/bin/* /usr/local/bin
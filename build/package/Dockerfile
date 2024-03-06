# Building image
FROM golang:alpine AS build

# Disable CGO
ENV CGO_ENABLED 0

# Get all remaining code
RUN mkdir -p /go/src/github.com/lerenn/asyncapi-codegen
COPY ./ /go/src/github.com/lerenn/asyncapi-codegen

# Set the workdir
WORKDIR /go/src/github.com/lerenn/asyncapi-codegen

# Build everything in cmd/
RUN --mount=type=cache,target=/root/.cache/go-build go install ./cmd/*

# Get final base image
FROM alpine

# Get binaries
COPY --from=build /go/bin/* /usr/local/bin
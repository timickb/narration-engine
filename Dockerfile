FROM golang:alpine AS builder
ENV CGO_ENABLED 0
ENV GOOS LINUX

RUN apk update --no-cache && apk add --no-cache tzdata

WORKDIR /build
ADD go.mod .
ADD go.sum .
RUN go mod download
COPY . .
FROM golang:1.17.2 AS builder

WORKDIR /build

ADD . .

RUN GO111MODULE=on \
  CGO_ENABLED=0 \
  go build -o connector ./cmd/connector

FROM alpine

WORKDIR /app

COPY --from=build /build/connector .

ENTRYPOINT ["/app/connector"]

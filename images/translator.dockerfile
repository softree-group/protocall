FROM golang:1.17 AS build

WORKDIR /build

ADD . .

RUN GO111MODULE=on \
  CGO_ENABLED=0 \
  go build -o translator ./cmd/translator

FROM alpine:3.14

WORKDIR /app

COPY --from=build /build/translator .

ENTRYPOINT ["/app/translator"]

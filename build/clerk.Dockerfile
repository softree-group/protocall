FROM golang:1.17 AS build
ENV CGO_ENABLED 0
ENV GO111MODULE on
ARG BUILD_REF
ARG IMAGE
ADD . /build

WORKDIR /build
RUN go build -o ${IMAGE} -ldflags "-X main.build=${BUILD_REF}" ./cmd/${IMAGE}

FROM alpine:3.14
ARG BUILD_DATE
ARG BUILD_REF
ARG IMAGE
COPY --from=build /build/${IMAGE} /usr/bin
WORKDIR /${IMAGE}
ENTRYPOINT ["clerk"]

LABEL org.opencontainers.image.created="${BUILD_DATE}" \
  org.opencontainers.image.title="${IMAGE}" \
  org.opencontainers.image.source="https://github.com/softree-group/protocall-connector" \
  org.opencontainers.image.revision="${BUILD_REF}" \
  org.opencontainers.image.vendor="Softex"

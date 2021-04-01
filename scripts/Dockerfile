ARG ALPINE_VERSION=3.10
ARG GOLANG_VERSION=1.15.5

FROM golang:${GOLANG_VERSION} AS builder
WORKDIR /go/src/github.com/vliubezny/gnotify
COPY . .
RUN make linux

FROM alpine:${ALPINE_VERSION}
RUN apk update && apk add --no-cache ca-certificates
WORKDIR /app
COPY static /app/static
COPY --from=builder /go/src/github.com/vliubezny/gnotify/build/gnotify-linux-amd64 /app/gnotify
ENTRYPOINT [ "/app/gnotify" ]
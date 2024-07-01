FROM golang:1.22-alpine as builder

RUN set -ex; \
    apk update; \
    apk add --no-cache \
      make \
      bash \
      nodejs \
      npm

WORKDIR /app
COPY . .

RUN set -ex; \
    make install-vue; \
    make build-vue

RUN set -ex; \
    go env -w GO111MODULE=on; \
    make statik; \
    CGO_ENABLED=0 make gocron2-node; \
    CGO_ENABLED=0 make gocron2

FROM sstc/headful-chromium:latest

RUN set -ex; \
	apt-get update; \
	DEBIAN_FRONTEND=noninteractive apt install --no-install-recommends -y \
		tini

WORKDIR /app
COPY --from=builder /app/bin bin
EXPOSE 5920
EXPOSE 5921
ENTRYPOINT ["tini", "--"]

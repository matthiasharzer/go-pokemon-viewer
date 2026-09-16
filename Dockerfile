FROM node:26.7.0-trixie-slim AS build-ui

WORKDIR /app

COPY . .

WORKDIR /app/ui

RUN npm ci && \
		npm run build

FROM golang:1.26.5-alpine3.24 AS build-server

ARG version=unknown

RUN apk update && \
		apk add git

WORKDIR /go/src

COPY go.mod go.sum ./
RUN go mod download && \
		go mod verify

COPY . .
COPY --from=build-ui /app/ui/public ui/public

RUN module_path=$(go list -m) && \
	go build \
		-o /go/bin/go-pokemon-viewer \
		-ldflags "-X ${module_path}/cmd/version.version=$version" \
		.

FROM alpine:3.24

RUN addgroup -S app && adduser -S -G app app

COPY --from=build-server /go/bin/go-pokemon-viewer /usr/local/bin/go-pokemon-viewer

WORKDIR /var/lib/go-pokemon-viewer
RUN chown app:app /var/lib/go-pokemon-viewer

USER app

ENTRYPOINT ["/usr/local/bin/go-pokemon-viewer"]

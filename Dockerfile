# web
FROM node:16.14.0-alpine as build-js
ENV PATH /app/node_modules/.bin:$PATH

WORKDIR /app

COPY ./web/package*.json ./
RUN npm ci --silent

COPY ./web ./
RUN npm run build

# server/cli
FROM golang:1.18 AS build-server

RUN echo 'deb [trusted=yes] https://repo.goreleaser.com/apt/ /' > /etc/apt/sources.list.d/goreleaser.list && \
  apt-get -y update && \
  apt-get -y install goreleaser-pro

WORKDIR /go/src

COPY ./server ./server
COPY ./cli ./cli
COPY ./.goreleaser.yaml .
RUN goreleaser build --single-target --rm-dist --snapshot

FROM alpine AS release

# Enable machine-id on alpine-linux (https://gitlab.alpinelinux.org/alpine/aports/-/issues/8761)
RUN apk add dbus

WORKDIR /app

COPY --from=build-server /go/src/dist/tracetest-server ./
COPY --from=build-server /go/src/dist/tracetest ./

COPY --from=build-server /go/src/server/migrations/ ./migrations/

COPY --from=build-js /app/build /app/html

EXPOSE 11633/tcp

ENTRYPOINT ["/app/tracetest-server"]

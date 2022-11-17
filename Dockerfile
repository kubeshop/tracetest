FROM node:16.14.0-alpine as build-js
WORKDIR /app
ENV PATH /app/node_modules/.bin:$PATH
COPY ./web/package.json ./
COPY ./web/package-lock.json ./
RUN npm ci --silent
COPY ./web ./
RUN npm run build

FROM golang:1.18-alpine AS build-server
WORKDIR /go/src

ARG ANALYTICS_BE_KEY
ARG ANALYTICS_FE_KEY
ARG VERSION
ARG TRACETEST_ENV
ARG POKE_API

RUN apk add --update make

COPY ./server/go.mod ./server/go.sum ./
RUN go mod download
COPY ./server ./
RUN make build

FROM golang:1.18-alpine AS build-cli
WORKDIR /go/src

ARG ANALYTICS_BE_KEY
ARG VERSION
ARG TRACETEST_ENV

RUN apk add --update make git jq

COPY ./server/go.mod ./server/go.sum ./server/
RUN cd server && go mod download
COPY ./server ./server

COPY ./cli/go.mod ./cli/go.sum ./cli/
RUN cd cli && go mod download

COPY ./cli ./cli

COPY .goreleaser.yaml .
RUN cd cli && make build

FROM alpine AS release
# Enable machine-id on alpine-linux (https://gitlab.alpinelinux.org/alpine/aports/-/issues/8761)
RUN apk add dbus
WORKDIR /app
COPY --from=build-server /go/src/tracetest-server ./
COPY --from=build-server /go/src/migrations/ ./migrations/
COPY --from=build-cli /go/src/cli/dist/tracetest ./
COPY --from=build-js /app/build /app/html
EXPOSE 11633/tcp
ENTRYPOINT ["/app/tracetest-server"]

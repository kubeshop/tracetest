FROM node:16.14.0-alpine as build-js
WORKDIR /app
ENV PATH /app/node_modules/.bin:$PATH
COPY ./web/package.json ./
COPY ./web/package-lock.json ./
RUN npm ci --silent
COPY ./web ./
RUN npm run build

FROM golang:1.18-alpine AS build-go
WORKDIR /go/src

COPY ./server/go.mod ./server/go.sum ./
RUN go mod download
COPY ./server ./
RUN go build -mod=readonly -o tracetest-server .

FROM alpine AS release
# Enable machine-id on alpine-linux (https://gitlab.alpinelinux.org/alpine/aports/-/issues/8761)
RUN apk add dbus
WORKDIR /app
COPY --from=build-go /go/src/tracetest-server ./
COPY --from=build-go /go/src/migrations/ ./migrations/
COPY --from=build-js /app/build /app/html
EXPOSE 8080/tcp
ENTRYPOINT ["/app/tracetest-server"]

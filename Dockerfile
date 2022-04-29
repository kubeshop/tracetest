FROM node:16.14.0-alpine as build-js
WORKDIR /app
ENV PATH /app/node_modules/.bin:$PATH
COPY ./web/package.json ./
COPY ./web/package-lock.json ./
RUN npm ci --silent
COPY ./web ./
RUN npm run build

FROM golang:1.17 AS build-go
WORKDIR /go/src

COPY ./server/go.mod ./server/go.sum ./
RUN go mod download
COPY ./server ./
RUN go build -o tracetest-server .

FROM ubuntu AS release
WORKDIR /app
COPY --from=build-go /go/src/tracetest-server ./
COPY --from=build-js /app/build /app/html
EXPOSE 8080/tcp
ENTRYPOINT ["/app/tracetest-server"]

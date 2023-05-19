FROM goreleaser/goreleaser:v1.11.2 AS build-cli
WORKDIR /app

RUN apk add --update jq make

COPY ./.goreleaser.yaml ./

COPY ./server/go.mod ./server/go.sum ./server/
RUN cd server && go mod download
COPY ./server ./server

COPY ./cli/go.mod ./cli/go.sum ./cli/
RUN cd cli && go mod download

COPY ./cli ./cli
RUN cd ./cli && make build

FROM golang:1.18-alpine

RUN apk --update add bash jq curl

WORKDIR /app
COPY --from=build-cli /app/cli/dist/tracetest /app/cli/tracetest
COPY ./testing/server-tracetesting ./tracetesting

WORKDIR /app/tracetesting
CMD ["/bin/sh", "/app/tracetesting/run.bash"]

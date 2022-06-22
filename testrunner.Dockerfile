FROM golang:1.18-alpine AS build-cli
WORKDIR /go/src

RUN apk add --update make

COPY ./cli/go.mod ./cli/go.sum ./
RUN go mod download
COPY ./cli ./
RUN make build

FROM alpine

RUN apk add bash jq curl

WORKDIR /app
COPY --from=build-cli /go/src/tracetest /app/cli/tracetest
COPY ./tracetesting ./tracetesting

WORKDIR /app/tracetesting
CMD ["/bin/bash", "/app/tracetesting/run.bash"]

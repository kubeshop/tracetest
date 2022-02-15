FROM golang:1.17 AS build
WORKDIR /go/src

COPY go.mod go.sum ./
RUN go mod download
COPY go ./go
COPY main.go ./main.go

RUN go build -o openapi .

FROM ubuntu AS release
WORKDIR /app
COPY --from=build /go/src/openapi ./
EXPOSE 8080/tcp
ENTRYPOINT ["/app/openapi"]

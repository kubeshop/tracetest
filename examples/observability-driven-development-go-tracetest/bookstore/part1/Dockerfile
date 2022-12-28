FROM golang:1.19

ARG SERVICE

WORKDIR /app/${SERVICE}

COPY ./${SERVICE}/go.* /app/${SERVICE}
RUN go mod download

COPY ./${SERVICE}/* /app/${SERVICE}
RUN go build -o /app/server .

ENTRYPOINT [ "/app/server" ]

FROM golang:1.21-alpine as builder

WORKDIR /app

ENV CGO_ENABLED 0
RUN go install go.k6.io/xk6/cmd/xk6@latest

RUN xk6 build v0.50.0 --with github.com/kubeshop/xk6-tracetest

FROM alpine

COPY --from=builder /app/k6 /bin/
COPY ./import-pokemon.js .
ENV XK6_TRACETEST_API_TOKEN=your-api-token

ENTRYPOINT [ "k6", "run", "/import-pokemon.js", "-o", "xk6-tracetest"]

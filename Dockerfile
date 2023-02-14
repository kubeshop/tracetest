FROM alpine AS release

WORKDIR /app

COPY ./tracetest-server /app/tracetest-server
COPY ./tracetest /app/tracetest

COPY ./html ./html

EXPOSE 11633/tcp

ENTRYPOINT ["/app/tracetest-server", "serve"]

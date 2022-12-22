FROM alpine AS release
# Enable machine-id on alpine-linux (https://gitlab.alpinelinux.org/alpine/aports/-/issues/8761)
RUN apk add dbus

WORKDIR /app

COPY ./tracetest-server /app/tracetest-server
COPY ./tracetest /app/tracetest

COPY ./web/build ./html
COPY ./server/migrations/ /app/migrations/


EXPOSE 11633/tcp
ENTRYPOINT ["/app/tracetest-server"]

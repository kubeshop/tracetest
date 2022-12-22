FROM alpine AS release
# Enable machine-id on alpine-linux (https://gitlab.alpinelinux.org/alpine/aports/-/issues/8761)
RUN apk add dbus

ARG OS_ARCH=linux_amd64

WORKDIR /app

COPY ./dist/${OS_ARCH}/tracetest-server ./
COPY ./dist/${OS_ARCH}/tracetest ./

COPY ./web/build ./html

COPY ./server/migrations/ /app/migrations/

EXPOSE 11633/tcp
ENTRYPOINT ["/app/tracetest-server"]

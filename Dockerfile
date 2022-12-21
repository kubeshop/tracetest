FROM alpine AS release
# Enable machine-id on alpine-linux (https://gitlab.alpinelinux.org/alpine/aports/-/issues/8761)
RUN apk add dbus

ARG OS=linux
ARG ARCH=amd64_v1

COPY dist/${OS}/server_${OS}_${ARCH}/tracetest-server /app/
COPY dist/${OS}/cli_${OS}_${ARCH}/tracetest /app/

COPY web/build /app/html

COPY server/migrations/ /app/migrations/

EXPOSE 11633/tcp
ENTRYPOINT ["/app/tracetest-server"]

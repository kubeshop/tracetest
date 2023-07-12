FROM alpine

RUN apk --update add bash jq curl

WORKDIR /app
COPY ./dist/tracetest /app/tracetest
COPY ./testing/server-tracetesting ./tracetesting

WORKDIR /app/tracetesting
CMD ["/bin/sh", "/app/tracetesting/run.bash"]

ARG TRACETEST_VERSION=latest

FROM kubeshop/tracetest:${TRACETEST_VERSION}
WORKDIR /app

ENV TRACETEST_API_KEY ""

ENTRYPOINT [ "tracetest", "start", "--api-key", "$TRACETEST_API_KEY" ]


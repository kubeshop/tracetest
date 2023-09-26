FROM kubeshop/tracetest:v0.13.9
WORKDIR /app

ENV TRACETEST_API_KEY ""

# ENTRYPOINT [ "ls" ]
ENTRYPOINT [ "tracetest", "start", "--api-key", "$TRACETEST_API_KEY" ]


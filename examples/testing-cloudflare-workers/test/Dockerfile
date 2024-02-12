FROM alpine
WORKDIR /app
RUN apk --update add bash jq curl
RUN curl -L https://raw.githubusercontent.com/kubeshop/tracetest/main/install-cli.sh | bash
WORKDIR /app/test/
ENTRYPOINT ["/bin/bash", "/app/test/run.bash"]

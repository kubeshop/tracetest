FROM ubuntu

RUN apt-get update && apt-get -y install \
  curl \
  && curl -L https://raw.githubusercontent.com/kubeshop/tracetest/main/install-cli.sh | sh

WORKDIR /app
COPY ./tracetesting ./tracetesting

WORKDIR /app/tracetesting
CMD ["/bin/bash", "/app/tracetesting/run.bash"]

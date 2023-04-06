# Provisioning server

Tracetest allows a server to be provisioned the first time a new Tracetest server is installed and launched. Provisioning sets certain resources in the server to the specified values, allowing you to configure the server. It is convenient in a CI/CD flow where you want to launch a server with a specified configuration. 

The server is provisioned by specifying a series of YAML snippets which will configure various resources. Each snippet is separated with the YAML separator, `---`.

Currently, the following resources can be provisioned: 
- DataStore
- PollingProfile
- Config
- Demo

For Docker-based installs, the provisioning file is placed in the `./tracetest/tracetest-provisioning.yaml` file by default when you run the `tracetest server install` command and select the `Using Docker Compose` option. The first time you start Tracetest with a `docker compose -f tracetest/docker-compose.yaml  up -d` command, the server will use the contents of this file to provision the server. To provision differently, you would alter the contents of the `tracetest-provisioning.yaml` file before launching Tracetest in Docker.

This is an example of a `tracetest-provisioning.yaml` file:

```yaml
---
type: DataStore
spec:
  name: otlp
  type: otlp
  otlp:
    type: otlp
---
type: Config
spec:
  analyticsEnabled: true
---
type: PollingProfile
spec:
  name: Custom Profile
  strategy: periodic
  default: true
  periodic:
    timeout: 2m
    retryDelay: 3s
---
type: Demo
spec:
  name: pokeshop
  type: pokeshop
  enabled: true
  pokeshop:
    httpEndpoint: http://demo-api:8081
    grpcEndpoint: demo-api:8082
```

Alternatively, we support setting an environment variable called `TRACETEST_PROVISIONING` to provision the server when it is first started. Base64 encode the provisioning YAML you want to utilize and set the `TRACETEST_PROVISIONING` environment variable with the result. The Tracetest server will then provision based on the Base64 encoded provisioning data in this environment variable the first time it is launched.


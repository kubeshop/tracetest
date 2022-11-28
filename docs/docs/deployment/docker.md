# Docker Deployment

You can run Tracetest locally using Docker Compose. This setup is great for a development environment. In this form, Tracetest runs in parallel to your Dockerized application,
allowing you to interact with your app and its traces, create and run tests over them, etc.

![Installer using docker compose](../img/installer/1_docker-compose_0.7.0.png)

**Tools required (installed if missing)**:
- Docker
- Docker Compose

**Requirements**:
- Jaeger or other compatible backend. If missing, the installer will help you configure one.
- OpenTelemetry Collector. If missing, the installer will help you configure one.
- A `docker-compose.yaml` (configurable) file in the project directory. If missing, the installer will create an empty file.

**Optionals**:
- [PokeShop demo app](https://github.com/kubeshop/pokeshop/)

**Result**:
- `tracetest/` directory (configurable) with a `docker-compose.yaml` and other config files.
- [Jaeger](https://www.jaegertracing.io/) instance, if selected.
- [OpenTelemetry Collector](https://opentelemetry.io/docs/collector/), if selected.
- [PokeShop demo app](https://github.com/kubeshop/pokeshop/), if selected.


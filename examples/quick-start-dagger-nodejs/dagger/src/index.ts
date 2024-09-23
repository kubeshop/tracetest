/**
 * A generated module for QuickStartDaggerNodejs functions
 *
 * This module has been generated via dagger init and serves as a reference to
 * basic module structure as you get started with Dagger.
 *
 * Two functions have been pre-created. You can modify, delete, or add to them,
 * as needed. They demonstrate usage of arguments and return types using simple
 * echo and grep commands. The functions can be called from the dagger CLI or
 * from one of the SDKs.
 *
 * The first line in this comment block is a short description line and the
 * rest is a long description with more detail on the module's purpose or usage,
 * if appropriate. All modules should have a short description.
 */
import { dag, Container, Directory, object, func, Service, File, ClientTracetestOpts, ContainerWithExecOpts, Secret } from "@dagger.io/dagger"
import { Context } from "../sdk/context/context"
import * as dotenv from "dotenv";
dotenv.config({ path: __dirname+'../../env' });

const otelCollectorConfig = `
receivers:
  otlp:
    protocols:
      grpc:
      http:

processors:
  batch:
    timeout: 100ms

  # Data sources: traces
  probabilistic_sampler:
    hash_seed: 22
    sampling_percentage: 100

exporters:
  otlp/tracetestagent:
    endpoint: tracetest-agent:4317
    tls:
      insecure: true

service:
  pipelines:
    traces/tracetestagent:
      receivers: [otlp]
      processors: [probabilistic_sampler, batch]
      exporters: [otlp/tracetestagent]
`

const tracetestTestDefinition = `
type: Test
spec:
  id: W656Q0c4g
  name: Test API
  description: akadlkasjdf
  trigger:
    type: http
    httpRequest:
      url: http://app:8080
      method: GET
      headers:
      - key: Content-Type
        value: application/json
  specs:
  - selector: span[tracetest.span.type="http" name="GET /" http.target="/" http.method="GET"]
    assertions:
    - attr:http.status_code  =  200
    - attr:tracetest.span.duration  <  500ms
`

@object()
class QuickStartDaggerNodejs {
  @func()
  build(source: Directory): Container {
    const nodeCache = dag.cacheVolume("node")
    return dag
      .container()
      .from("node:21-slim")
      .withDirectory("/usr/src/app", source)
      .withMountedCache("/usr/src/app/node_modules", nodeCache)
      .withWorkdir("/usr/src/app")
      .withExec(["npm", "install"])
  }

  @func()
  app(source: Directory): Service {
    return this.build(source)
      .withExec(["npm", "run", "with-grpc-tracer"])
      .withEnvVariable('OTEL_EXPORTER_OTLP_TRACES_ENDPOINT', process.env.OTEL_EXPORTER_OTLP_TRACES_ENDPOINT)
      .withExposedPort(8080)
      .asService()
  }

  @func()
  tracetestAgent(): Service {
    return dag.container()
      .from('kubeshop/tracetest:v1.6.1')
      .withEnvVariable('TRACETEST_TOKEN', process.env.TRACETEST_TOKEN)
      .withEnvVariable('TRACETEST_ENVIRONMENT_ID', process.env.TRACETEST_ENVIRONMENT_ID)
      .withExposedPort(4317)
      .asService()
  }

  @func()
  tracetestCli(source: Directory): Promise<string> {
    this.app(source)
    this.tracetestAgent()

    return dag.container()
      .from('alpine')
      .withWorkdir("/app")
      .withExec(["apk", "--update", "add", "bash", "jq", "curl"])
      .withExec(["curl", "-L", `https://raw.githubusercontent.com/kubeshop/tracetest/main/install-cli.sh | bash -s -- ${process.env.TRACETEST_IMAGE_VERSION}`])
      .withExec(["configure", "--token", process.env.TRACETEST_TOKEN, "--environment", process.env.TRACETEST_ENVIRONMENT_ID])
      .withNewFile("test.yaml", tracetestTestDefinition)
      .withExec(["run", "test", "--file", "test.yaml"])
      .stdout()
  }

  // @func()
  // otelCollector(): Service {
  //   return dag.container()
  //     .from('otel/opentelemetry-collector:0.100.0')
  //     .withNewFile("/etc/otelcol/config.yaml", otelCollectorConfig)
  //     .withExec(["/otelcol", "--config", "/etc/otelcol/config.yaml"])
  //     .withExposedPort(4317)
  //     .withServiceBinding(
  //       "tracetest-agent",
  //       dag.tracetest().agent().asService()
  //     )
  //     .asService()
  // }

//   @func()
//   tracetest(source: Directory, token: Secret, environment: string): Promise<string> {
//     this.app(source)
//     this.otelCollector()

//     return dag.tracetest({
//       apiKey: token,
//       environment: environment,
//     })
//       .cli()
//       .withNewFile("test.yaml", tracetestTestDefinition)
//       .withExec(["run", "test", "--file", "test.yaml"])
//       .stdout()
//   }
// }

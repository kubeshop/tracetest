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

@object()
class QuickStartDaggerNodejs {

  apiKey: any
  environment: any
  organization: any

  constructor(apiKey, environment, organization) {
    this.apiKey = apiKey
    this.environment = environment
    this.organization = organization
  }

  @func()
  build(source: Directory): Container {
    const nodeCache = dag.cacheVolume("node")
    return dag
      .container()
      .from("node:20-slim")
      .withDirectory("/usr/src/app", source)
      .withMountedCache("/usr/src/app/node_modules", nodeCache)
      .withWorkdir("/usr/src/app")
      .withExec(["npm", "install"])
  }

  @func()
  app(
    source: Directory
  ): Service {
    return this.build(source)
      .withExec(["npm", "run", "with-grpc-tracer"])
      .withEnvVariable('OTEL_EXPORTER_OTLP_TRACES_ENDPOINT', 'http://tracetest-agent:4317')
      .withExposedPort(8080)
      .withServiceBinding(
        'tracetest-agent',
        dag
          .tracetest(
            this.apiKey,
            this.environment,
            this.organization
          )
          .agent()
          .asService()
      )
      .asService()
  }

  @func()
  tracetest(
    source: Directory
  ): Promise<string> {

    const appSvc = this.app(source)
    const	appSvcEndpoint = appSvc.endpoint()

    const tracetestTestDefinition = `
      type: Test
      spec:
        id: W656Q0c4g
        name: Test API
        description: akadlkasjdf
        trigger:
          type: http
          httpRequest:
            url: ${appSvcEndpoint}
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

    return dag
      .tracetest(
        this.apiKey,
        this.environment,
        this.organization
      )
      .cli()
      .withNewFile("test.yaml", tracetestTestDefinition)
      .withExec(["run", "test", "--file", "test.yaml"], { "useEntrypoint": true })
      .stdout()
  }
}

# CI/CD Automation

This section contains a general overview of running Tracetest in CI/CD pipelines.

You can find guides for:

- [GitHub Actions](./github-actions-pipeline)
- [Testkube](./testkube-pipeline)
- [Tekton](./tekton-pipeline)

:::note
If you want to see more examples with other CI/CD tools, let us know by [opening an issue in GitHub](https://github.com/kubeshop/tracetest/issues/new/choose)!
:::

Tracetest is designed to work with all CI/CD platforms and automation tools. To enable Tracetest to run in CI/CD environments, make sure to [install the Tracetest CLI](../getting-started/installation.mdx) and configure it to access your [Tracetest server](../configuration/server.md).

You can also directly execute the Tracetest CLI from a Docker image rather than installing the CLI on your local machine. This can be convenient when you wish to execute the CLI in a CI/CD environment.

**How to Use**:

```sh
docker run --rm -it -v$(pwd):$(pwd) -w $(pwd) --entrypoint tracetest kubeshop/tracetest:latest -s http://host.docker.internal:11633/ test run  --definition <file-path> --wait-for-result
```

To read more about integrating Tracetest with CI/CD tools, check out tutorials in our blog:

- [Integrating Tracetest with GitHub Actions in a CI pipeline](https://kubeshop.io/blog/integrating-tracetest-with-github-actions-in-a-ci-pipeline)




# Opening Tracetest

This page showcases opening the Tracetest Web UI regardless if you used the Tracetest CLI, Docker, Kubernetes, or Helm to install Tracetest Server.

Once you've installed Tracetest, as explained in the [installation guide](./installation.mdx), you access the Tracetest Web UI on [`http://localhost:11633`](http://localhost:11633). Here's what will greet you after a fresh install.

![Landing page Tracetest](https://res.cloudinary.com/djwdcmwdz/image/upload/v1688474565/docs/screely-1688474539641_kbhvvc.png)

By following the [installation guide](./installation.mdx) your Tracetest instance will have a `demo` Pokeshop app installed that generates distributed traces when triggered.

## Creating Trace-based Tests

You can create tests in two ways:

- Visually, in the Web UI
- Programmatically, in YAML

## Creating Visual Trace-based Tests

This guide will show how to create end-to-end and integration tests in less than 5 minutes via the Web UI.

:::note
To view the in-depth guide on creating tests visually, [check out this docs page](../web-ui/creating-tests.md).
:::

### Create

On the top right, click the **Create** button and select **Create New Test** in the drop down.

![Create a new test](https://res.cloudinary.com/djwdcmwdz/image/upload/v1688475179/docs/screely-1688475174365_ckq3cn.png)

Select an **HTTP Request** as the **test trigger**, and choose the **Pokeshop - Import** example.

![Select Pokeshop example](https://res.cloudinary.com/djwdcmwdz/image/upload/v1688475514/docs/screely-1688475510090_r6hqmx.png)

This will populate a sample API test against a POST endpoint in the Pokeshop app demo. Clicking **Create & Run** will save and trigger the test.

![API test against POST endpoint](https://res.cloudinary.com/djwdcmwdz/image/upload/v1688475680/docs/screely-1688475676524_vvtxsu.png)

:::info
Running a test against `localhost` will resolve as `127.0.0.1` inside the Tracetest container. To run tests against apps running on your local machine, add them to the same network and use service name mapping instead. Example: Instead of running an app on `localhost:8080`, add it to your Docker Compose file, connect it to the same network as your Tracetest service, and use `service-name:8080` in the URL field when creating an app.

You can reach services running on your local machine using:

- Linux (docker version < 20.10.0): `172.17.0.1:8080`
- MacOS (docker version >= 18.03) and Linux (docker version >= 20.10.0): `host.docker.internal:8080`
:::

### Trigger

You'll be redirected to the test page where you can see four tabs and depending on which one you select you'll get access to:

- Test trigger and results
- The entire distributed trace and trace analysis
- Test specification and assertions
- How to automate the test

:::note
To view the in-depth guide on test results, [check out this docs page](http://localhost:3000/web-ui/test-results)
:::


The **Trigger** tab shows how the test was triggered and what the API response was.

![test result](https://res.cloudinary.com/djwdcmwdz/image/upload/v1688476389/docs/screely-1688476384678_edcsgx.png)

### Trace

The **Trace** tab shows the entire distributed trace for debugging and a trace analysis score.

![distributed trace and trace analysis score](https://res.cloudinary.com/djwdcmwdz/image/upload/v1688476460/docs/screely-1688476455986_q24aa2.png)

### Test

The **Test** tab shows span attributes. Here you add test specs and assertions on attribute values. You also get a test snippets out-of-the-box for common test cases.

Here you see how to assert that all database spans return in less than 100ms.

![test specs](https://res.cloudinary.com/djwdcmwdz/image/upload/v1688476657/docs/screely-1688476653521_omxe4r.png)

### Automate

The **Automate** tab shows how to automate the test run with the Tracetest CLI and other automation options.

![automate](https://res.cloudinary.com/djwdcmwdz/image/upload/v1688476810/docs/screely-1688476801601_f4s0iy.png)




# Detailed Instructions on Installing Tracetest Using the CLI

Tracetest has a command line interface (CLI) with includes an  **install wizard** that helps with installing the Tracetest server into Docker or Kubernetes. The CLI can also be used run tests, download or upload tests, and many manage much of the capability of Tracetest.

## Installing the Tracetest Server via the CLI

Use the CLI's install wizard to install a Tracetest server locally using Docker Compose or to a local or remote Kubernetes cluster.

The wizard installs all the tools required to set up the desired environment and creates all the configurations, tailored to your case.

Every time we release a new version of Tracetest, we generate binaries for Linux, MacOS, and Windows. Supporting both amd64, and ARM64 architectures, in `tar.gz`, `deb`, `rpm` and `exe` formats.

You can find the latest version [here](https://github.com/kubeshop/tracetest/releases/latest).

### Linux/MacOS

Tracetest CLI can be installed automatically using the following script:

```sh
curl -L https://raw.githubusercontent.com/kubeshop/tracetest/main/install-cli.sh | bash
```

It works for systems with Homebrew, `apt-get`, `dpkg`, `yum`, `rpm` installed, and if no package manager is available, it will try to download the build and install it manually.

You can also manually install with one of the following methods.

#### Homebrew

```sh
brew install kubeshop/tracetest/tracetest
```

#### APT

```sh
# requirements for our deb repo to work
sudo apt-get update && sudo apt-get install -y apt-transport-https ca-certificates

# add repo
echo "deb [trusted=yes] https://apt.fury.io/tracetest/ /" | sudo tee /etc/apt/stores.list.d/fury.list

# update and install
sudo apt-get update
sudo apt-get install tracetest
```

#### YUM

```sh
# add repository
cat <<EOF | $SUDO tee /etc/yum.repos.d/tracetest.repo
[tracetest]
name=Tracetest
baseurl=https://yum.fury.io/tracetest/
enabled=1
gpgcheck=0
EOF

# install
sudo yum install tracetest --refresh
```

### Windows
Download one of the files from the latest tag, extract to your machine, and then [add the tracetest binary to your PATH variable](https://stackoverflow.com/a/41895179)

## Install a Tracetest Server for Development with the CLI

This guide will help you get Tracetest running using the Tracetest CLI.

### Prerequisites

:::info
Make sure you have [Docker](https://docs.docker.com/get-docker/) and [Docker Compose](https://docs.docker.com/compose/install/), or [Kubernetes](https://kubernetes.io/) installed.
:::

:::info
In this quick start, OpenTelemetry Collector is used to send traces directly to Tracetest. If you have an existing trace data source, [read here](../configuration/overview.md).
:::

### 1. Run the `server install` Command
Once you've installed the CLI you can install a Tracetest server by running:

```bash
tracetest server install
```

```text title="Expected output:"
████████ ██████   █████   ██████ ███████ ████████ ███████ ███████ ████████
   ██    ██   ██ ██   ██ ██      ██         ██    ██      ██         ██
   ██    ██████  ███████ ██      █████      ██    █████   ███████    ██
   ██    ██   ██ ██   ██ ██      ██         ██    ██           ██    ██
   ██    ██   ██ ██   ██  ██████ ███████    ██    ███████ ███████    ██




Hi! Welcome to the TraceTest server installer. I'll help you set up your TraceTest server by asking you a few questions
and configuring your system with all the requirements, so you can start TraceTesting right away!

To get more info about TraceTest, you can check our docs at https://kubeshop.github.io/tracetest/

If you have any issues, please let us know by creating an issue (https://github.com/kubeshop/tracetest/issues/new/choose)
or reach us on Discord https://discord.gg/6zupCZFQbe


How do you want to run TraceTest? [type to search]:
> Using Docker Compose
  Using Kubernetes
```

### 2. Select Docker Compose

This prompts you to select if you want to get either a Docker Compose or Kubernetes setup generated for you.

Select Docker Compose.

```text title="Expected output:"
...

-> Let's check if your system has everything we need

✔ docker already installed
✔ docker is ready
✔ docker compose already installed

-> Your system is ready! Now, let's configure TraceTest

Project's docker-compose file [docker-compose.yaml]:
```

### 3. Generate a Docker Compose File

Add the name of the Docker Compose file, if you have one, or just hit enter to proceed and add one.

```text title="Expected output:"
ERROR   File "docker-compose.yaml" does not exist. You need an existing docker-compose file.
Do you want me to create an empty docker-compose file? [Y/n]:
```

Now, hit enter again and this will generate a default `docker-compose.yaml` file.

### 4. Configure Jaeger as the Trace Data Source

Next up you'll be prompted to connect a trace data store, if you have one.

```text title="Expected output:"
Do you have a supported tracing backend you want to use? (Jaeger, Tempo, OpenSearch, SignalFX) [y/N]:
```

Write `Y` and hit enter, or just proceed if you do not have a trace data store configured.

```text title="Expected output:"
Do you want me to set up Jaeger? [Y/n]:
```

We can set up Jaeger for you if you wish. Proceed to use Jaeger as the trace data store. We'll generate a Jaeger config for you.

:::info
If you have an existing trace data source, [read this](../configuration/overview.md), or see how to use OpenTelemetry Collector [here](../configuration/connecting-to-data-stores/opentelemetry-collector).
:::

### 5. Configure OpenTelemetry Collector

Next, you'll be prompted about OpenTelemetry Collector. Proceed with `N`.

```text title="Expected output:"
Do you have an OpenTelemetry Collector? [y/N]:
```

Now we prompt you if you want us to generate an OpenTelemetry Collector config for you. Select `Y`.

```text title="Expected output:"
Do you want me to set up one? [Y/n]:
```

This will generate an OpenTelemetry Collector config for you.

:::info
Want to read in more detail how to use OpenTelemetry Collector? [Check this out](../configuration/connecting-to-data-stores/opentelemetry-collector).
:::

### (Optional) Enable a Demo App

In the next step, the CLI will ask if you want a demo app to try out your Tracetest config.

```text title="Expected output:"
Do you want to enable the PokeShop demo app? (https://github.com/kubeshop/pokeshop/) [Y/n]:
```

Proceed with `Y` again.

### 6. Select Output Directory for Tracetest Server Installation

Lastly, select the output directory where to store all files related to Tracetest. 

```text title="Expected output:"
Tracetest output directory [tracetest/]:
```

Proceeding will set it to `tracetest/` by default.

```text title="Expected output:"
-> Thanks! We are ready to install TraceTest now

 SUCCESS  Install successful!

To start tracetest:

	docker compose -f docker-compose.yaml -f tracetest/docker-compose.yaml up -d

Then, use your browser to navigate to:

  http://localhost:11633

Happy TraceTesting =)
```

:::note
This means all Tracetest-related config files, files for Docker Compose, and any config file for OpenTelemetry Collector will be located in a `./tracetest/` directory, with a `docker-compose.yaml` file located in the root directory where you ran the `tracetest server install` command.
:::

### 7. Start Docker Compose

```bash
docker compose -f docker-compose.yaml -f tracetest/docker-compose.yaml up -d
```

```bash title="Condensed expected output from the Tracetest container:"
Starting tracetest ...
...
2022/11/28 18:24:09 HTTP Server started
...
```

### 8. Open the Tracetest Web UI

Open [`http://localhost:11633`](http://localhost:11633) in your browser.

Create a [test](../web-ui/creating-tests.md).

:::info
Running a test against `localhost` will resolve as `127.0.0.1` inside the Tracetest container. To run tests against apps running on your local machine, add them to the same network and use service name mapping instead. Example: Instead of running an app on `localhost:8080`, add it to your Docker Compose file, connect it to the same network as your Tracetest service, and use `service-name:8080` in the URL field when creating an app.

You can reach services running on your local machine using:

- Linux (docker version < 20.10.0): `172.17.0.1:8080`
- MacOS (docker version >= 18.03) and Linux (docker version >= 20.10.0): `host.docker.internal:8080`
:::

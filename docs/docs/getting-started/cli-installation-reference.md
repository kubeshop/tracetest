# CLI Installation Reference

This page contains a reference of all options how to install Tracetest CLI.

## Detailed Instructions on Installing Tracetest Using the CLI

Tracetest has a command line interface (CLI) which includes an **install wizard** that helps to install the Tracetest server into Docker or Kubernetes. The CLI can also be used to run tests, download or upload tests, and manage much of the capability of Tracetest.

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

You can also manually install it with one of the following methods.

#### Homebrew

```sh
brew install kubeshop/tracetest/tracetest
```

#### APT

```sh
# requirements for our deb repo to work
sudo apt-get update && sudo apt-get install -y apt-transport-https ca-certificates

# add repo
echo "deb [trusted=yes] https://apt.fury.io/tracetest/ /" | sudo tee /etc/apt/sources.list.d/fury.list

# update and install
sudo apt-get update
sudo apt-get install tracetest
```

#### YUM

```sh
# add repository
cat <<EOF | sudo tee /etc/yum.repos.d/tracetest.repo
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

#### Chocolatey

```bash
choco source add --name=kubeshop_repo --source=https://chocolatey.kubeshop.io/chocolatey ; choco install tracetest
```

#### From source

Download one of the files from the latest tag, extract to your machine, and then [add the tracetest binary to your PATH variable](https://stackoverflow.com/a/41895179).

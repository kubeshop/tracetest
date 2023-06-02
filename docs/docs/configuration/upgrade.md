# Upgrade Tracetest Version

This page explains how to upgrade the version of your Tracetest Server and CLI.

If you've ever seen this error, you've come to the right page:

```text
✖️ Error: Version Mismatch
The CLI version and the server version are not compatible. To fix this, you'll need to make sure that both your CLI and server are using compatible versions.
We recommend upgrading both of them to the latest available version. Check out our documentation https://docs.tracetest.io/configuration/upgrade for simple instructions on how to upgrade.
Thank you for using Tracetest! We apologize for any inconvenience caused.
```

This means your Tracetest CLI and Server versions must be the same.

```sh
tracetest version
```

```text title="Expected output"
CLI: v0.11.9
Server: v0.11.9
✔️ Version match
```

## Upgrade Tracetest CLI Version

### Linux/MacOS

Run the Tracetest CLI install script to upgrade to the latest version of the CLI:

```sh
curl -L https://raw.githubusercontent.com/kubeshop/tracetest/main/install-cli.sh | bash
```

### Homebrew

```sh
brew upgrade
brew update
brew install kubeshop/tracetest/tracetest
```

### APT

```sh
sudo apt-get update
sudo apt-get install tracetest
```

### YUM

```sh
sudo yum update
sudo yum install tracetest --refresh
```

### Windows

```sh
choco source add --name=kubeshop_repo --source=https://chocolatey.kubeshop.io/chocolatey ; choco upgrade tracetest
```

## Upgrade Tracetest Server Version

Make sure to match the CLI version you have installed to the Server version.

```sh
kubeshop/tracetest:vX.X.X
```

If you are using version `v0.11.9` of the CLI, make sure to use the same version of the server.

```yaml
image: kubeshop/tracetest:v0.11.9
```

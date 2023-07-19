#!/bin/sh

version=$1

cmd_exists() {
  command -v $1 &> /dev/null
}

ensure_dependency_exist() {
  if ! cmd_exists $1; then
    echo "missing dependency: $1 is required to run this script"
    exit 2
  fi
}

get_latest_version() {
  if [[ -n "$version" ]]; then
    echo $version
    exit 0
  fi

  ensure_dependency_exist "curl"

  curl --silent "https://api.github.com/repos/kubeshop/tracetest/releases/latest" |
  grep '"tag_name":' |
  sed -E 's/.*"([^"]+)".*/\1/'
}


get_os() {
  os_name=`uname | tr '[:upper:]' '[:lower:]'`
  echo $os_name
}

get_arch() {
  arch=$(uname -p)
  if [ "$arch" = "unknown" ]; then
    arch=$(uname -m)
  fi
  case "$arch" in
    "x86_64")
      echo "amd64"
      ;;

    "arm"*|"aarch64")
      echo "arm64"
      ;;

    *)
      echo $arch
      ;;
  esac
}

get_download_link() {
  os=$(get_os)
  version=$(get_latest_version)
  arch=$(get_arch)
  raw_version=`echo $version | sed 's/v//'`
  pkg=$1


  echo "https://github.com/kubeshop/tracetest/releases/download/${version}/tracetest_${raw_version}_${os}_${arch}.${pkg}"
}

download_file() {
  file=$1
  path=$2

  echo "Downloading $file and saving to $path"

  curl -L "$file" --output "$path"
  echo "File downloaded and saved to $path"
}

install_tar() {
  ensure_dependency_exist "curl"
  ensure_dependency_exist "tar"

  download_link=$(get_download_link "tar.gz")
  file_path="/tmp/cli.tar.gz"
  download_file "$download_link" "$file_path"

  echo "Extracting file"
  tar -xf $file_path -C /tmp
  echo "Installing to /usr/local/bin/tracetest"
  $SUDO mv /tmp/tracetest /usr/local/bin/tracetest
  rm -f $file_path
}

install_dpkg() {
  download_link=$(get_download_link "deb")
  file_path="/tmp/cli.deb"
  download_file "$download_link" "$file_path"

  $SUDO dpkg -i $file_path
  rm -f $file_path
}

install_rpm() {
  download_link=$(get_download_link "rpm")
  file_path="/tmp/cli.rpm"
  download_file "$download_link" "$file_path"

  $SUDO rpm -i $file_path
  rm -f $file_path
}

install_apt() {
  $SUDO apt-get update
  $SUDO apt-get install -y apt-transport-https ca-certificates
  echo "deb [trusted=yes] https://apt.fury.io/tracetest/ /" | $SUDO tee /etc/apt/sources.list.d/fury.list
  $SUDO apt-get update
  $SUDO apt-get install -y tracetest
}

install_yum() {
  cat <<EOF | $SUDO tee /etc/yum.repos.d/tracetest.repo
[tracetest]
name=Tracetest
baseurl=https://yum.fury.io/tracetest/
enabled=1
gpgcheck=0
EOF
  $SUDO yum install -y tracetest --refresh
}

install_brew() {
  brew install kubeshop/tracetest/tracetest
}

run() {
  if [ ! -z "$version" ]; then
    echo "Installing version $version"
    install_tar
  fi
  ensure_dependency_exist "uname"
  if cmd_exists brew; then
    install_brew
  elif cmd_exists apt-get; then
    install_apt
  elif cmd_exists yum; then
    install_yum
  elif cmd_exists dpkg; then
    install_dpkg
  elif cmd_exists rpm; then
    install_rpm
  elif [ "$(get_arch)" == "unknown" ]; then
    echo "unknown system architecture. Try manual install. See https://kubeshop.github.io/tracetest/installing/#cli-installation"
    exit 1;
  elif [[ "$(get_os)" =~ ^(darwin|linux)$ ]]; then
    install_tar
  else
    echo 'OS not supported by this script. See https://kubeshop.github.io/tracetest/installing/#cli-installation'
    exit 1
  fi

  echo
  echo "Succesfull install!"
  echo
  echo "run 'tracetest --help' to see what you can do"
  echo
  echo "To setup a new server, run 'tracetest server install'"
}

SUDO=""
if [ `id -un` != "root" ]; then
  ensure_dependency_exist "sudo"
  SUDO="sudo"
fi

run

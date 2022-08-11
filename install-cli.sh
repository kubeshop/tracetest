#!/bin/bash

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
    arch=`uname -p`
    if [ "x86_64" == "$arch" ]; then
        echo "amd64"
    elif [ "arm" == "$arch" ]; then
        echo "arm64"
    else
        echo "$arch"
    fi
}

get_download_link() {
    os=$(get_os)
    version=$(get_latest_version)
    arch=$(get_arch)
    pkg=$4
    raw_version=`echo $version | sed 's/v//'`

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
  download_link=$(get_download_link "tar.gz")
  file_path="/tmp/cli.tar.gz"
  download_file "$download_link" "$file_path"

  tar -xvf $compressed_file_path -C /tmp
  $SUDO mv /tmp/tracetest /usr/local/bin/tracetest
}

install_dpkg() {
  download_link=$(get_download_link "deb")
  file_path="/tmp/cli.deb"
  download_file "$download_link" "$file_path"

  $SUDO dpkg -i $file_path
}

install_rpm() {
  download_link=$(get_download_link "rpm")
  file_path="/tmp/cli.rpm"
  download_file "$download_link" "$file_path"

  $SUDO rpm -i $file_path
}

install_apt() {
  $SUDO apt-get update
  $SUDO apt-get install -y apt-transport-https ca-certificates
  echo "deb [trusted=yes] https://apt.fury.io/tracetest/ /" | $SUDO tee /etc/apt/sources.list.d/fury.list
  $SUDO apt-get update
  $SUDO apt-get install -y tracetest
}

install_yum() {
  cat <<EOF | $SUDO tee /etc/yum.repos.d/fury.repo
[fury]
name=Tracetest
baseurl=https://yum.fury.io/tracetest/
enabled=1
gpgcheck=0
EOF
  $SUDO yum install -y tracetest --refresh
}

run() {
    ensure_dependency_exist "uname"
    os=$(get_os)
    if [ "$os" = "windows" ]; then
      echo 'OS not supported by this script. See https://kubeshop.github.io/tracetest/installing/#cli-installation'
      exit 1
    fi

    if cmd_exists apt-get; then
      install_apt
    elif cmd_exists yum; then
      install_yum
    elif cmd_exists dpkg; then
      install_dpkg
    elif cmd_exists rpm; then
      install_rpm
    else
      install_tar
    fi
}

SUDO=""
if [ `id -un` != "root" ]; then
  ensure_dependency_exist "sudo"
  SUDO="sudo"
fi

run

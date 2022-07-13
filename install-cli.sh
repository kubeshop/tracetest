#!/bin/bash

ensure_dependency_exist() {
    if ! command -v $1 &> /dev/null
    then
        echo "missing dependency: $1 is required to run this script"
        exit 2
    fi
}

ensure_required_dependencies_are_present() {
    ensure_dependency_exist "curl"
    ensure_dependency_exist "tar"
    ensure_dependency_exist "uname"
}

get_latest_version() {
    redirect_url=`curl -Ls -o /dev/null -w %{url_effective} https://github.com/kubeshop/tracetest/releases/latest`
    url_suffix="https://github.com/kubeshop/tracetest/releases/tag/"
    version=${redirect_url/#$url_suffix}

    echo $version
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
    os=$1
    arch=$2
    version=$3

    echo "https://github.com/kubeshop/tracetest/releases/download/$version/tracetest-$version-$os-$arch.tar.gz"
}

download_file() {
    file=$1
    path=$2

    echo "Downloading $file and saving to $path"

    curl -L "$file" --output "$path"
    echo "File downloaded and saved to $path"
}

install_cli() {
    compressed_file_path=$1

    tar -xvf $compressed_file_path
    sudo cp tracetest /usr/local/bin/tracetest
}

run() {
    ensure_required_dependencies_are_present

    latest_version=`get_latest_version`
    os=`get_os`
    arch=`get_arch`
    download_link=`get_download_link $os $arch $latest_version`
    file_path="/tmp/cli.tar.gz"

    echo "Downloading Tracetest CLI $latest_version for $os $arch"
    echo "URL link: $download_link"

    download_file "$download_link" "$file_path"
    install_cli "$file_path"
}

run

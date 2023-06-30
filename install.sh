#!/usr/bin/env bash

VERSION_URL="https://raw.githubusercontent.com/cristianoliveira/ergo/master/.version"
DOWNLOAD_URL="https://github.com/cristianoliveira/ergo/releases/download"
DEST_FOLDER="/usr/local/bin"
PROGNAME=`basename "$0"`
VERSION=${1:-$(wget -q -O - "$VERSION_URL")}

die () {
echo "$PROGNAME: [FATAL] $1" >&2; exit ${2:-1}  ;
}

getplatform(){
    platform=$(uname -m)
    case $platform in
    x86_64)
        echo "linux-amd64"
        ;;
    *)
        die "$platform is not yet supported"
        ;;
    esac
}

install(){
    echo "Using /tmp to store downloaded file"
    cd /tmp
    local platform=$(getplatform)
    echo "Downloading version $VERSION from repo"
    wget -q "$DOWNLOAD_URL/$VERSION/ergo-$VERSION-$platform.tar.gz"
    [ $? -ne 0 ] && die "unable to download package"

    echo "Extracting package"
    tar -xf ergo-$VERSION-$platform.tar.gz
    [ $? -ne 0 ] && die "unable to extract ergo from package"

    echo "Copying ergo to $DEST_FOLDER. May require sudo password."
    if [ -w $DEST_FOLDER ]; then
        cp ergo $DEST_FOLDER
    else
        sudo cp ergo $DEST_FOLDER
    fi
    [ $? -ne 0 ] && die "unable to copy ergo to destination folder"

    echo "Application was successfully installed."
    echo "For uninstalling execute: rm ${DEST_FOLDER}ergo"
}

show_help(){
    echo "Usage: $PROGNAME [-d destination_directory]"
}

main(){
    [ -z "$VERSION" ] && die "Unable to get the version information to install"

    install
}

while getopts "h?d:" opt; do
    case "$opt" in
    h|\?)
        show_help
        exit 0
        ;;
    d)  DEST_FOLDER=$OPTARG
        ;;
    esac
done

main

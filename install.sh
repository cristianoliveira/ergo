#!/usr/bin/env bash

VERSION_URL="https://raw.githubusercontent.com/cristianoliveira/ergo/master/.version"
DOWNLOAD_URL="https://github.com/cristianoliveira/ergo/releases/download"
declare -A SUPPORTED_PLATFORMS
SUPPORTED_PLATFORMS=(["x86_64"]="linux-amd64", ["amd64"]="freebsd-amd64")
DEST_FOLDER="/usr/local/bin"
PROGNAME=`basename "$0"`

function die () { 
echo "$PROGNAME: [FATAL] $1" >&2; exit ${2:-1}  ; 
}

function install(){
    echo "Using /tmp to store downloaded file"
    cd /tmp

    echo "Downloading version $latest_version from repo"
    
    wget -q "$DOWNLOAD_URL/$latest_version/ergo-$latest_version-${SUPPORTED_PLATFORMS[$platform]}.tar.gz"
    [ $? -ne 0 ] && die "unable to download package"

    echo "Extracting package"    
    tar -xf ergo-$latest_version-${SUPPORTED_PLATFORMS[$platform]}.tar.gz 
    [ $? -ne 0 ] && die "unable to extract ergo from package"

    echo "Copying ergo to $DEST_FOLDER. May require sudo password."
    if [ -w $DEST_FOLDER ]; then
        cp ergo $DEST_FOLDER
    else
        sudo cp ergo $DEST_FOLDER
    fi
    [ $? -ne 0 ] && die "unable to copy ergo to destination folder"

    echo "Application was successfully installed."
    echo "For uninstalling execute: rm $DEST_FOLDER/ergo"
}

function show_help(){
    echo "Usage: $PROGNAME [-d destination_directory]"
}

function main(){
    platform=$(uname -m)
    [ ${SUPPORTED_PLATFORMS[$platform]} ] || die "$platform is not yet supported"
    
    latest_version=$(wget -q -O - "$VERSION_URL")
    [ $? -ne 0 ] && die "unable to retrieve latest version information"

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

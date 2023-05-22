#!/usr/bin/env bash
set -ex

make all

cp bin/ergo ergo

# sanity check the file type
file ergo

# release tarball will look like 'go-ergo-v1.2.3-x86_64-unknown-linux-gnu.tar.gz'
ARTIFACT="ergo-${RELEASE_TAG:?"Missing release tag"}-${TARGET}.tar.gz"
tar czf "$ARTIFACT" ergo

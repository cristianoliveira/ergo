# `before_deploy` phase: here we package the build artifacts

set -ex

# create a "staging" directory
mkdir staging

cp bin/ergo staging

cd staging

# release tarball will look like 'rust-everywhere-v1.2.3-x86_64-unknown-linux-gnu.tar.gz'
tar czf ../${PROJECT_NAME}-${TRAVIS_TAG}-${TARGET}.tar.gz *

# `script` phase: you usually build, test and generate docs in this phase

set -ex

make all
make coverage

# sanity check the file type
file bin/ergo

# `script` phase: you usually build, test and generate docs in this phase

set -ex

go build -o bin/ergo

# sanity check the file type
file bin/ergo

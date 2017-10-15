#!/usr/bin/env bash

echo "After success: Coverage metrics report..."

set -e
echo "" > coverage.txt

for d in $(go list ../... | grep -v vendor); do
    go test -race -coverprofile=profile.out -covermode=atomic $d
    if [ -f profile.out ]; then
        cat profile.out >> coverage.txt
        rm profile.out
    fi
done

bash < $(curl -s https://codecov.io/bash) -t ${CODECOV_TOKEN}

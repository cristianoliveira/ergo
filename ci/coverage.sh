#!/usr/bin/env bash

echo "After success: Coverage metrics report..."

set -e
echo "" > coverage.txt

for d in $(go list ../... | grep -v -e '\(examples\|vendor\)'); do
    echo "testing: $d"
    go test -coverprofile=profile.out -covermode=atomic -v $d
    if [ -f profile.out ]; then
        cat profile.out >> coverage.txt
        rm profile.out
    fi
done

bash ./ci/codecov.sh -t ${CODECOV_TOKEN}

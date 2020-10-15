#!/usr/bin/env bash

set -eu

cd "$(dirname "$0")/.."

tempdir=$(mktemp -d)
go build -o "$tempdir/tengo-tester" ./cmd/tengo-tester
export PATH="$tempdir:$PATH"
cd examples
# command -v tengo-tester
go test -race -covermode=atomic ./...
rm -R "$tempdir"

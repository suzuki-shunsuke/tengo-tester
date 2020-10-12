#!/usr/bin/env bash

set -eu
set -o pipefail

ee() {
  echo "+ $*"
  eval "$@"
}

cd "$(dirname "$0")/.."

repo_name=${1:-}
if [ -z "$repo_name" ]; then
  echo "the repository name is required" >&2
  exit 1
fi

ee cc-test-reporter before-build

mkdir -p .code-climate

for d in $(go list ./...); do
  echo "$d"
  profile=.code-climate/$d/profile.txt
  coverage=.code-climate/$d/coverage.json
  ee mkdir -p "$(dirname "$profile")" "$(dirname "$coverage")"
  ee go test -race -coverprofile="$profile" -covermode=atomic "$d"
  if [ "$(wc -l < "$profile")" -eq 1 ]; then
    continue
  fi
  ee cc-test-reporter format-coverage -t gocov -p "github.com/suzuki-shunsuke/${repo_name}" -o "$coverage" "$profile"
done

result=.code-climate/codeclimate.total.json
if [ 0 -eq "$(find .code-climate -name coverage.json | wc -l)" ]; then
  echo "no test found" >&2
  exit 0
fi
# shellcheck disable=SC2046
ee cc-test-reporter sum-coverage $(find .code-climate -name coverage.json) -o "$result"
ee cc-test-reporter upload-coverage -i "$result"

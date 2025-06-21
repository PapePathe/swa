#! /bin/bash

set -euxo pipefail

go build -x -o tmp/swa main.go

./tmp/swa "$@"

#! /bin/bash

set -euxo pipefail

go build -o tmp/swa main.go

./tmp/swa "$@"

#! /bin/bash

set -euxo pipefail

go test ./... -count=1
go build -o tmp/swa main.go

cp ./tmp/swa ../tests
cd ../tests

go test ./... -count=1 "$@"

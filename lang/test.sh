#! /bin/bash

set -euxo pipefail

if [ -f "start.s" ]; then
    echo "deleting assembly files"
    rm -f start.s
fi

go test ./...
go build -o tmp/swa main.go

cp ./tmp/swa ../tests
cd ../tests

go test ./... -count=1 "$@"

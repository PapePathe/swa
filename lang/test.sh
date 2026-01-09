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

if compgen -G "*.s"; then
    rm -rf *.s
fi

if compgen -G "*.ll"; then
    rm -rf *.ll
fi

if compgen -G "*.exe"; then
    rm -rf *.exe
fi

if compgen -G "*.o"; then
    rm -rf *.o
fi

go test ./... -count=1 "$@"

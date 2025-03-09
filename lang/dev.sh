#! /bin/sh

go build -o tmp/swa main.go

./tmp/swa "$@"

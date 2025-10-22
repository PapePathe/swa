#! /bin/bash

set -euxo pipefail

go build -o swahili main.go

fpm -t deb
fpm -t apk

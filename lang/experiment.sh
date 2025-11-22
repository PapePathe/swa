#! /bin/bash

set -euxo pipefail

go build -o tmp/swa main.go

./tmp/swa compile -s test.swa && cat start.ll && ./start.exe

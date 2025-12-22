#! /bin/bash

set -euxo pipefail

go build -o tmp/swa main.go

if [ "$1" == "parse" ]; then
    echo "🚀 Parsing..."
    ./tmp/swa parse -s test.swa
elif [ "$1" == "compile" ]; then
    echo "🛑 Compiling..."
    ./tmp/swa compile -e -s test.swa && ./start.exe
else
    echo "❓ Unknown command: $1. Use 'compile' or 'parse'."
    exit 1
fi

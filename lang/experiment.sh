#! /bin/bash

set -euxo pipefail

go build -o tmp/swa main.go

if [ "$1" == "parse" ]; then
    echo "🚀 Parsing..."
    ./tmp/swa parse -s test.swa
elif [ "$1" == "tokenize" ]; then
    echo "🚀 Tokenizing..."
    ./tmp/swa tokenize -s test.swa
elif [ "$1" == "compile" ]; then
    echo "🛑 Compiling..."
    ./tmp/swa compile -s test.swa && cat start.ll && ./start.exe
else
    echo "❓ Unknown command: $1. Use 'compile' or 'parse'."
    exit 1
fi

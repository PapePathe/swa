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
    ./tmp/swa compile -s test.swa && ./start.exe
elif [ "$1" == "compile-exp" ]; then
    echo "🛑 Compiling Experimental version..."
    SWA_DEBUG=yes ./tmp/swa compile -e -s test.swa && ./start.exe

else
    echo "❓ Unknown command: $1. Use 'compile' or 'parse'."
    exit 1
fi

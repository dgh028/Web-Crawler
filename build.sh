#!/usr/bin/bash

rm -rf ./bin
rm -rf ./log
rm -rf ./output
mkdir log
mkdir output

go mod tidy
# go run main.go
# cd main
# go build -o ../bin/mini_spider
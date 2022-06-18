#!/bin/sh
cd proto
protoc --go_out=. *.proto
cd ..
go mod tidy
go run .
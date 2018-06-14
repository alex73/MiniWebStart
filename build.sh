#!/bin/sh

rm -rf dist
mkdir dist
GOOS=windows GOARCH=386   go build -o dist/mini-win32.exe -ldflags "-s -w" src/*.go
GOOS=windows GOARCH=amd64 go build -o dist/mini-win64.exe -ldflags "-s -w" src/*.go
GOOS=darwin  GOARCH=amd64 go build -o dist/mini-mac64     -ldflags "-s -w" src/*.go
GOOS=linux   GOARCH=amd64 go build -o dist/mini-linux64   -ldflags "-s -w" src/*.go
GOOS=linux   GOARCH=386   go build -o dist/mini-linux32   -ldflags "-s -w" src/*.go

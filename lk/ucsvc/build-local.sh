#!/bin/bash
set -e
set -v

# Since this does not exit with a non 0 exit code, we cannot stop
go vet .

# Since this does not exit with a non 0 exit code, we cannot stop
golint .

# Ensure Godeps symbolic link exists
[ ! -d Godeps ] && ln -s ../Godeps Godeps

# Run go build via godep
godep go build .

# Run all tests excluding integration tests via godep
godep go test -v

# Run all tests including integration tests via godep, this one is a little odd as a service instance needs to be running
godep go test -tags "test integration" -v

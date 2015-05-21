#!/bin/bash
set -e
set -v

# Since this does not exit with a non 0 exit code, we cannot stop
go vet .

# Since this does not exit with a non 0 exit code, we cannot stop
golint .

# Ensure Godeps symbolic link exists
[ ! -d Godeps ] && ln -s ../Godeps

# Run go build via godep
godep go build .

# Run all tests via godep - None yet
godep go test -v

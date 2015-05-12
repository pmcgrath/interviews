#!/bin/bash

echo Running go vet .
go vet .


echo Running golint .
golint .


echo Ensuring Godeps symbolic link exists
[ ! -d Godeps ] && ln -s ../Godeps Godeps


echo Running go build via godep
godep go build .
[ $? -ne 0 ] && exit


echo All including integration tests
echo Running go test -tags "test integration" -v
go test -tags "test integration" -v
[ $? -ne 0 ] && exit


echo All excluding integration tests
echo Running go test -v
go test -v
[ $? -ne 0 ] && exit

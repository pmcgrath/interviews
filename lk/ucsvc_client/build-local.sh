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

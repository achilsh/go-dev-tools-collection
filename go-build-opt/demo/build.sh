#! /usr/bin/env bash 
pkgPath="github.com/achilsh/go-dev-tools-collection/go-build-opt"
flags="-X $pkgPath.Version=v0.1 -X '$pkgPath.GoVersion=$(go version)' -X '$pkgPath.BuildTime=`date +"%Y-%m-%d %H:%m:%S"`' -X $pkgPath.GitCommit=`git rev-parse HEAD`"

go build -ldflags "$flags" -o version_demo
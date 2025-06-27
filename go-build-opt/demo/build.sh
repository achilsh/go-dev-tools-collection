#! /usr/bin/env bash 
pkgPath="github.com/achilsh/go-dev-tools-collection/go-build-opt"
flags="-X $pkgPath.Version=v0.1 -X '$path.GoVersion=$(go version)' -X '$path.BuildTime=`date +"%Y-%m-%d %H:%m:%S"`' -X $path.GitCommit=`git rev-parse HEAD`"

go build -ldflags "$flags" -o version_demo
#!/bin/bash
CURDIR=$(cd $(dirname $0); pwd)
BinaryName=hz-new-proto
echo "$CURDIR/bin/${BinaryName}"
exec $CURDIR/bin/${BinaryName}
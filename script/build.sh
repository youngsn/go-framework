#!/bin/bash

# Use to build project.
# As added -a pram when compiling, that means compile will force to rebuild all packages.

ROOT=$(cd `dirname $0`; cd ..; pwd)
NAME=${ROOT##*/}

# build sh must use vendor package
export GOPATH=$ROOT:$ROOT/vendor

go build -o $ROOT/bin/$NAME -v -x $ROOT/src/$NAME/main.go

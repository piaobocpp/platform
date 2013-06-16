#!/usr/bin/env bash

if [ ! -f install.sh ]; then
    echo 'install must be run within its container folder' 1>&2
    exit 1
fi

CURDIR=`pwd`
OLDGOPATH="$GOPATH"
export GOPATH="$CURDIR"

gofmt -w src

gofmt -tabs=false -tabwidth=4 -w src

go install -v Platform

echo 'finished'

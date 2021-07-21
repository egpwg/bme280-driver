#!/bin/bash
env GOOS=linux GOARCH=arm64 go build -ldflags "-w -s"

if [ ! -d bin ];then
    mkdir bin
fi
mv -f demo bin
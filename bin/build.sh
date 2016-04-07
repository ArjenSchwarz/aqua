#!/bin/bash
set -ex

GOOS=linux GOARCH=amd64 go build -a -ldflags '-s' -o lambda/aqua

zip -j lambda/aqua.zip lambda/index.js lambda/aqua

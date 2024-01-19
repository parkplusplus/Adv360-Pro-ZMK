#!/bin/bash
set -ex

git checkout -f ./config
git pull
go run ./fixkeyboard/fixkeyboard.go
make

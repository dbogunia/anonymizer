#!/bin/bash
cd /go/src/anonymizer
git pull
go run main.go $1
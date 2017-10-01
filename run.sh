#!/bin/bash
goimports -w .
go fmt github.com/algon-320/...
go build *.go

./kide $@

cp kide ~/KIDE/

@echo off
cls
gofmt -d -w -s .
go build
CuteASM
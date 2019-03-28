@echo off
protoc -I=protoc --go_out=protoc2 protoc2.proto
protoc -I=protoc --go_out=protoc3 protoc3.proto
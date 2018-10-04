#!/usr/bin/env bash
#/bin/bash
##生成proto文件
protoc --proto_path=. --micro_out=. --go_out=. ./proto/saywo/saywo.proto
##生成二进制代码
#go build main.go
##构建docker镜像
docker build -t micro-service/order-service:latest .




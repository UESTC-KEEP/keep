#!/bin/bash
rm -rf prome-exporter
go build ../prome-exporter.go
#mv prome-exporter ../prome-exporter
# docker build -f /path/to/a/Dockerfile .
# 构建镜像
docker build  -f ../Dockerfile -t keep-prome-exporter .
# 删除运行失败的镜像 若失败可注释
#docker rm  `docker ps -a | grep  'keep-cloud' | awk '{print $1}'`
# 打标签传habor
docker tag keep-prome-exporter:latest 172.17.15.242/library/keep-prome-exporter:latest
docker push 172.17.15.242/library/keep-prome-exporter:latest
#   删除none镜像
docker rmi `docker images | grep  '<none>' | awk '{print $3}'`
#docker run --name keep-cloud -p 20350:20350 -p 8888:8888  `docker images | grep  'keep-cloud' | awk '{print $3}'`
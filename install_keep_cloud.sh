#!/bin/bash
rm -rf cloudagent
# docker build -f /path/to/a/Dockerfile .
docker build -f Dockerfile_cloud -t keep-cloud .
docker rm  `docker ps -a | grep  'keep-cloud' | awk '{print $1}'`
docker run --name keep-cloud -p 20350:20350 -p 8888:8888  `docker images | grep  'keep-cloud' | awk '{print $3}'`

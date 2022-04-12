# 本镜像用于在集群环境测试promethues exporter 不做为生产镜像的一部分
# 打包之前务必更改包名为main
FROM golang:1.17

##工作目录作目录
#WORKDIR $GOPATH/src/keep

COPY . .

EXPOSE 8080
ENTRYPOINT ["./prome-exporter"]

# docker build -t keep-cloud .
# docker run --name keep-cloud -p 20350:20350 -p 8888:8888  `docker images | grep  'keep-cloud' | awk '{print $3}'`
#   一条命令实现停用并删除所有容器：
#   docker stop $(docker ps -q) & docker rm $(docker ps -aq)
#   删除none镜像
#   docker rmi `docker images | grep  '<none>' | awk '{print $3}'`
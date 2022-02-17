# 脚本仅仅用于第一次安装minikube环境  注意脚本将会卸载现有docker 如需保留请注释卸载部分（kubernetes版本指定1.21.0）
# 虚拟机至少2核 2G
# 后续启动minikube测试环境 执行下面一条命令即可：
# minikube start  --force --driver=docker --image-mirror-country='cn'   --registry-mirror=https://registry.docker-cn.com  --kubernetes-version=v1.21.0
echo  "安装kubectl 参考链接：https://kubernetes.io/docs/tasks/tools/install-kubectl-linux/"
curl -LO https://dl.k8s.io/release/v1.21.0/bin/linux/amd64/kubectl
sudo install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl
kubectl version --client
echo "安装minikube 参考链接：https://minikube.sigs.k8s.io/docs/start/?spm=a2c6h.12873639.0.0.ab202043wymz7U "
curl -LO https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64
sudo install minikube-linux-amd64 /usr/local/bin/minikube
echo "安装docker"
#sudo yum remove docker \
#                  docker-client \
#                  docker-client-latest \
#                  docker-common \
#                  docker-latest \
#                  docker-latest-logrotate \
#                  docker-logrotate \
#                  docker-engine
#sudo yum install -y yum-utils \
#  device-mapper-persistent-data \
#  lvm2
#sudo yum-config-manager \
#    --add-repo \
#    http://mirrors.aliyun.com/docker-ce/linux/centos/docker-ce.repo
#sudo yum install docker-ce docker-ce-cli containerd.io

#sudo systemctl start docker
#sudo systemctl enable docker


echo "配置minikube..."
minikube start  --force --driver=docker --image-mirror-country='cn'   --registry-mirror=https://registry.docker-cn.com  --kubernetes-version=v1.21.0
minikube node add
minikube addons enable metrics-server
minikube addons enable dashboard
minikube dashboard
kubectl get nodes



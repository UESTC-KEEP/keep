# 设置go代理
go env -w GO111MODULE=on
go env -w GOPROXY=https://goproxy.cn,direct
# 执行生成client脚本 -kubeconfig=/home/et/.kube/config -alsologtostderr=true
$GOPATH/src/k8s.io/code-generator/generate-groups.sh all keep/cloud/pkg/client keep/cloud/pkg/apis/keepedge equalnode:v1alpha1 tenantresourcequota:v1alpha1

$GOPATH/src/k8s.io/code-generator/generate-groups.sh all keep/cloud/pkg/client keep/cloud/pkg/apis/keepedge tenantresourcequota:v1alpha1

$GOPATH/src/k8s.io/code-generator/generate-groups.sh all keep/cloud/pkg/client/kubeedgeClient/devices keep/cloud/pkg/apis/kubeedge/devices/v1alpha2 kubeedge:v1alpha2
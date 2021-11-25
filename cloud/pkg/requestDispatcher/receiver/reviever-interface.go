package receiver

import (
	commonTypes "keep/pkg/apis/compoenentconfig/keep/v1alpha1/common"
	"net/http"
)

type ReceiverInterface interface {
	// StartReceiverAndListen 启动reciever监听
	/*
		传入参数：port代表所需要监听的端口集合
		返回值：error:表示整体启动过程的错误
			   map[int]error:如果上述error不是nil 这在map中返回对应port的错误信息
	*/
	StartReceiverAndListen(port []int) (error, map[int]error)
	//ResolveWebsocketMessage 解析收到的包
	/*
		监听到连接的websocket传递的数据之后需要对请求进行解码  解析出消息格式struct
	*/
	ResolveWebsocketMessage(*http.Request) (*commonTypes.WebSocketMessage, error)
	//AuthorizationRequest 对发起服务请求的服务进行校验是否具有合法的证书
	/*
		传入参数：msg:keep格式的消息
		返回值：对消息进行鉴权之后的结果和错误信息
	*/
	AuthorizationRequest(*commonTypes.WebSocketMessage) (bool, error)

	// SentMessageIntoBeehiveModule 将数据请求数据按照msg中的router信息进行转发
	/*
		传入参数：moduleName：存放在beehive注册时使用的模块名
	*/
	SentMessageIntoBeehiveModule(message *commonTypes.WebSocketMessage, moudleName string) error
}

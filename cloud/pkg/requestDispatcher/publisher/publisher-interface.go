package publisher

import (
	"github.com/gorilla/websocket"
	commonTypes "keep/pkg/apis/compoenentconfig/keep/v1alpha1/common"
	"time"
)

type PublisherInterface interface {
	// ConnectToEdgeReceiver 尝试与指定的边缘节点建立通道
	/*
		输入参数：connerction： 边缘设备主动链接receiver形成的websocket通道
				edgePort: 边缘暴露的端口
				timeOut: 容忍的建立连接时间
	*/
	ConnectToEdgeReceiver(connerction websocket.Conn, edgePort int, timeOut time.Duration) error
	// SendMsgToEdge 向边缘设备发送消息
	SendMsgToEdge(connerction websocket.Conn, timeOut time.Duration, msg commonTypes.WebSocketMessage) error
}

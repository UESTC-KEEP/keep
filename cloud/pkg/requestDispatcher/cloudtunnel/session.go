package cloudtunnel

import (
	"fmt"
	"github.com/UESTC-KEEP/keep/cloud/pkg/requestDispatcher/Router"
	"github.com/UESTC-KEEP/keep/pkg/stream"
	"github.com/UESTC-KEEP/keep/pkg/util/core/model"
	logger "github.com/UESTC-KEEP/keep/pkg/util/loggerv1.0.1"

	"github.com/gorilla/websocket"
)

type session struct {
	//表示与边缘端的一个连接id
	sessionID string
	//云与边缘之间的websocket隧道
	tunnel stream.SafeWriteTunneler
	//通道是否关闭
	tunnelClosed bool
}

func (s *session) writeMessageToTunnel(m *model.Message) error {
	return s.tunnel.WriteMessage([]*model.Message{m})
}

func (s *session) Close() {
	err := s.tunnel.Close()
	if err != nil {
		logger.Error("close tunnel failed:", err)
	}
	s.tunnelClosed = true
}

func (s *session) String() string {
	return fmt.Sprintf("Tunnel session [%v]", s.sessionID)
}

func (s *session) Serve() {
	defer s.Close()

	for {
		t, r, err := s.tunnel.NextReader()
		if err != nil {
			logger.Error(err)
			return
		}
		if t != websocket.TextMessage {
			logger.Error(err)
			return
		}
		messageList, err := stream.ReadMessageFromTunnel(r)
		if err != nil {
			logger.Error("read message from tunnel error: ", s.String(), err)
			return
		} else {
			//group := message.Router.Group
			//beehiveContext.SendToGroup(group, *message)
			// Router.MessageDispatcher(message)
			for _, message := range messageList {
				Router.RevMsgChan <- message
			}
		}

	}

}

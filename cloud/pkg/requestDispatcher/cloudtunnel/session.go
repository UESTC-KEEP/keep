package cloudtunnel

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/wonderivan/logger"
	"keep/pkg/stream"
	beehiveContext "keep/pkg/util/core/context"
	"keep/pkg/util/core/model"
)

type Session struct {
	//表示与边缘端的一个连接id
	sessionID string
	//云与边缘之间的websocket隧道
	tunnel stream.SafeWriteTunneler
	//通道是否关闭
	tunnelClosed bool
}

func (s *Session) WriteMessageToTunnel(m *model.Message) error {
	return s.tunnel.WriteMessage(m)
}

func (s *Session) Close() {
	err := s.tunnel.Close()
	if err != nil {
		logger.Error("close tunnel failed:", err)
	}
	s.tunnelClosed = true
}

func (s *Session) String() string {
	return fmt.Sprintf("Tunnel session [%v]", s.sessionID)
}

func (s *Session) Serve() {
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
		message, err := stream.ReadMessageFromTunnel(r)
		if err != nil {
			logger.Error("read message from tunnel: ", s.String(), err)
			return
		} else if message.Header.ParentID != "" {
			beehiveContext.SendResp(*message)
		} else {
			group := message.Router.Group
			beehiveContext.SendToGroup(group, *message)
		}

	}
}

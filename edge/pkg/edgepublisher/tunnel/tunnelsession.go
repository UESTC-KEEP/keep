package edgetunnel

import (
	"github.com/gorilla/websocket"
	"github.com/wonderivan/logger"
	"keep/pkg/stream"
	beehiveContext "keep/pkg/util/core/context"
	"sync"
	"time"
)

type TunnelSession struct {
	Tunnel    stream.SafeWriteTunneler
	closeLock sync.Mutex
	closed    bool
}

func NewTunnelSession(c *websocket.Conn) *TunnelSession {
	return &TunnelSession{
		closeLock: sync.Mutex{},
		Tunnel:    stream.NewDefaultTunnel(c),
	}
}

func (t *TunnelSession) Close() {
	t.closeLock.Lock()
	defer t.closeLock.Unlock()
	if !t.closed {
		err := t.Tunnel.Close()
		if err != nil {
			logger.Error("close tunnel error: ", err)
		}
	}
	t.closed = true
}

func (t *TunnelSession) startPing(reconnectChan chan struct{}) {
	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()

	for {
		select {
		case <-beehiveContext.Done():
			logger.Warn("EdgeTunnel startPing stop")
			return
		case <-ticker.C:
			err := t.Tunnel.WriteControl(websocket.PingMessage, []byte("ping you"), time.Now().Add(time.Second))
			if err != nil {
				logger.Error("ping error: ", err)
				reconnectChan <- struct{}{}
				return
			}
		}
	}
}

func (t *TunnelSession) routeToEdge(reconnectChan chan struct{}) {
	for {
		select {
		case <-beehiveContext.Done():
			logger.Warn("EdgeTunnel RouteToEdge stop")
			return
		default:
		}

		_, r, err := t.Tunnel.NextReader()
		if err != nil {
			logger.Error("Read messsage error: ", err)
			reconnectChan <- struct{}{}
			return
		}

		mess, err := stream.ReadMessageFromTunnel(r)
		if err != nil {
			logger.Error("Get tunnel message error: ", err)
			reconnectChan <- struct{}{}
			return
		}

		//如果是对某条消息的响应消息
		if mess.Header.ParentID != "" {
			beehiveContext.SendResp(*mess)
		} else {
			beehiveContext.SendToGroup(mess.Router.Group, *mess)
		}
	}
}

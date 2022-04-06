package edgetunnel

import (
	"github.com/UESTC-KEEP/keep/constants/edge"
	"github.com/UESTC-KEEP/keep/pkg/stream"
	beehiveContext "github.com/UESTC-KEEP/keep/pkg/util/core/context"
	"github.com/UESTC-KEEP/keep/pkg/util/loggerv1.0.1"
	"github.com/gorilla/websocket"
	"sync"
	"time"
)

type tunnelSession struct {
	Tunnel    stream.SafeWriteTunneler
	closeLock sync.Mutex
	closed    bool
}

func NewTunnelSession(c *websocket.Conn) *tunnelSession {
	return &tunnelSession{
		closeLock: sync.Mutex{},
		Tunnel:    stream.NewDefaultTunnel(c),
	}
}

func (t *tunnelSession) Close() {
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

func (t *tunnelSession) startPing(reconnectChan chan struct{}) {
	ticker := time.NewTicker(edge.DefaultEdgeHeartBeat * time.Second)
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

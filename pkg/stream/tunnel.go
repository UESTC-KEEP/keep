/*
Copyright 2020 The KubeEdge Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package stream

import (
	"bufio"
	"encoding/json"
	"github.com/wonderivan/logger"
	"io"
	"io/ioutil"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"keep/pkg/util/core/model"
)

type SafeWriteTunneler interface {
	WriteMessage(message *model.Message) error
	WriteControl(messageType int, data []byte, deadline time.Time) error
	NextReader() (messageType int, r io.Reader, err error)
	io.Closer
}

type DefaultTunnel struct {
	lock *sync.Mutex
	con  *websocket.Conn
}

func (t *DefaultTunnel) WriteControl(messageType int, data []byte, deadline time.Time) (e error) {
	t.lock.Lock()
	e = t.con.WriteControl(messageType, data, deadline)
	t.lock.Unlock()
	return
}

func (t *DefaultTunnel) Close() error {
	return t.con.Close()
}

func (t *DefaultTunnel) NextReader() (messageType int, r io.Reader, err error) {
	return t.con.NextReader()
}

func (t *DefaultTunnel) WriteMessage(m *model.Message) (e error) {
	t.lock.Lock()
	bytes, e := json.Marshal(m)
	e = t.con.WriteMessage(websocket.TextMessage, bytes)
	t.lock.Unlock()
	return
}

func NewDefaultTunnel(con *websocket.Conn) *DefaultTunnel {
	return &DefaultTunnel{
		lock: &sync.Mutex{},
		con:  con,
	}
}

var _ SafeWriteTunneler = &DefaultTunnel{}

func ReadMessageFromTunnel(r io.Reader) (msg *model.Message, err error) {
	buf := bufio.NewReader(r)
	data, err := ioutil.ReadAll(buf)
	if err != nil {
		return nil, err
	}
	logger.Info("receive from tunnel: ", string(data))
	msg = &model.Message{}
	err = json.Unmarshal(data, msg)
	return
}

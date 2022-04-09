package Router

import (
	beehiveContext "github.com/UESTC-KEEP/keep/pkg/util/core/context"
	"github.com/UESTC-KEEP/keep/pkg/util/core/model"
	"testing"
	"time"
)

func TestRouter_MessageDispatcher(t *testing.T) {

	msg := &model.Message{}
	msg.Content = "hello kafka"
	msg.Router.Group = "/log"

	msg1 := &model.Message{}
	msg1.Content = "hello kafka1"
	msg1.Router.Group = "/add"

	go MessageRouter()

	RevMsgChan <- msg

	time.Sleep(3 * time.Second)
	RevMsgChan <- msg1

}

func TestRouter_MessageDispatcher1(t *testing.T) {

	msg := &model.Message{}
	msg.Content = "hello kafka"
	msg.Router.Group = "/log"

	msg1 := &model.Message{}
	msg1.Content = "hello kafka1"
	msg1.Router.Group = "/log"

	go MessageRouter()

	RevMsgChan <- msg

	time.Sleep(3 * time.Second)
	RevMsgChan <- msg1

}

func TestRouter_SendToEdge2(t *testing.T) {
	msg := &model.Message{}
	msg.Content = "hello edge!!!"
	msg.Router.Group = "/log"

	beehiveContext.AddModule("router")
	beehiveContext.Send("router", *msg)

	time.Sleep(3 * time.Second)

	SendToEdge()
	time.Sleep(3 * time.Second)

	// if msg, err := beehiveContext.Receive("router"); err != nil {
	// 	fmt.Println("err:", err)
	// } else {
	// 	fmt.Println("msg :", msg.Content)
	// }

}

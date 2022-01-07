package manager

import (
	"github.com/UESTC-KEEP/keep/pkg/util/loggerv1.0.1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
)

type Manager interface {
	Events() chan watch.Event
}

type CommonResourceEventHandler struct {
	events chan watch.Event
}

func (c *CommonResourceEventHandler) obj2Event(t watch.EventType, obj interface{}) {
	eventObj, ok := obj.(runtime.Object)
	if !ok {
		logger.Error("传入类型不是 runtime.Object 无法转换....")
		return
	}
	c.events <- watch.Event{Type: t, Object: eventObj}
}

// OnAdd 处理对象新增事务
func (c *CommonResourceEventHandler) OnAdd(obj interface{}) {
	c.obj2Event(watch.Added, obj)
}

// OnUpdate 处理对象更新事务
func (c *CommonResourceEventHandler) OnUpdate(oldObj, newObj interface{}) {
	c.obj2Event(watch.Modified, newObj)
}

// OnDelete 处理对象删除事务
func (c *CommonResourceEventHandler) OnDelete(obj interface{}) {
	c.obj2Event(watch.Deleted, obj)
}

// NewCommonResourceEventHandler 新建handler对象
func NewCommonResourceEventHandler(events chan watch.Event) *CommonResourceEventHandler {
	return &CommonResourceEventHandler{events: events}
}

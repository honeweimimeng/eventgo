package eventgo

import (
	"github.com/honeweimimeng/eventgo/driver"
	"github.com/honeweimimeng/eventgo/driver/event"
	"testing"
	"time"
)

type TestEvent struct {
	id      int
	name    string
	channel driver.Channel
}

func (r *TestEvent) Name() string {
	return r.name
}

func (r *TestEvent) Id() int {
	return r.id
}

func (r *TestEvent) Channel() driver.Channel {
	return r.channel
}

func Events(id int) []event.Proto {
	return []event.Proto{
		&TestEvent{id: id, name: "WRITE"},
		&TestEvent{id: id, name: "WRITE"},
		&TestEvent{id: id, name: "CLOSE"},
	}
}

func HappenEvent(id int, name string, msg any) *TestEvent {
	e := &TestEvent{id: id, name: name}
	e.channel = &event.Channel{EventProto: e, MsgInstance: msg}
	return e
}

func TestNormalSocketEvent(t *testing.T) {
	var boot event.BootStrap
	boot.EventLoop().Handle(Events(0), func(ch driver.Channel) {
		println("===>")
		time.Sleep(2 * time.Second)
		println("===>finish")
	}).Trigger(func(ch chan []event.Proto) {
		ch <- []event.Proto{HappenEvent(0, "WRITE", "write msg")}
		ch <- []event.Proto{HappenEvent(0, "WRITE", "write msg")}
	}).Boot().StartUp()
	time.Sleep(200 * time.Second)
}

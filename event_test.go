package eventgo

import (
	"github.com/honeweimimeng/eventgo/driver"
	"github.com/honeweimimeng/eventgo/driver/event"
	"testing"
	"time"
)

type TestEvent struct {
	Id_      int
	Name_    string
	Channel_ driver.Channel
}

type TestEvent2 struct {
	TestEvent
}

func (r *TestEvent) Name() string {
	return r.Name_
}

func (r *TestEvent) Id() int {
	return r.Id_
}

func (r *TestEvent) Channel() driver.Channel {
	return r.Channel_
}

func Events(id int) []event.Proto {
	return []event.Proto{
		&TestEvent{Id_: id, Name_: "WRITE"},
		&TestEvent{Id_: id, Name_: "WRITE"},
		&TestEvent{Id_: id, Name_: "CLOSE"},
	}
}

func Events2(id int) []event.Proto {
	e := &TestEvent2{}
	e.Id_ = id
	e.Name_ = "T2WRITE"
	return []event.Proto{e}
}

func HappenEvent(id int, name string, msg any) *TestEvent {
	e := &TestEvent{Id_: id, Name_: name}
	e.Channel_ = &event.Channel{EventProto: e, MsgInstance: msg}
	return e
}

func HappenEvent2(id int, name string, msg any) *TestEvent2 {
	e := &TestEvent2{}
	e.Id_ = id
	e.Name_ = name
	e.Channel_ = &event.Channel{EventProto: e, MsgInstance: msg}
	return e
}

func TestNormalSocketEvent(t *testing.T) {
	var boot event.BootStrap
	boot.EventLoop().Handle(Events(0), func(ch driver.Channel) {
		println("===>eventLoop0")
	}).Trigger(func(ch chan []event.Proto) {
		ch <- []event.Proto{HappenEvent(0, "WRITE", "write msg")}
		time.Sleep(2 * time.Second)
		ch <- []event.Proto{HappenEvent(0, "WRITE", "write msg")}
	}).EventLoop().Handle(Events2(0), func(ch driver.Channel) {
		println("===>eventLoop1")
	}).Trigger(func(ch chan []event.Proto) {
		time.Sleep(2 * time.Second)
		ch <- []event.Proto{HappenEvent2(0, "T2WRITE", "write msg")}
		ch <- []event.Proto{HappenEvent(0, "WRITE", "write msg")}
	}).StartUp()
	time.Sleep(200 * time.Second)
}

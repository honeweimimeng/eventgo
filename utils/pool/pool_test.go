package pool

import (
	"container/heap"
	"context"
	"github.com/sirupsen/logrus"
	"litecluster/utils"
	"testing"
	"time"
)

type Comparable interface {
	Val() interface{}
	CompareTo(other interface{}) int
}

func GetComparablePipe() *ComparablePipe {
	arr := make([]Comparable, 0)
	r := &ComparablePipe{heap: arr}
	return r
}

type ComparablePipe struct {
	heap []Comparable
}

func (t *ComparablePipe) PushTask(task Task) {
	t.Push(task)
}

func (t *ComparablePipe) PopTask() Task {
	ro := t.Pop()
	if ro == nil {
		return nil
	}
	return ro.(Task)
}

func (t *ComparablePipe) Push(task any) {
	c, ok := task.(Comparable)
	if !ok {
		panic("task must type of Comparable")
	}
	t.heap = append(t.heap, c)
	heap.Init(t)
}

func (t *ComparablePipe) Pop() any {
	var v Comparable = nil
	var r Task = nil
	if t.Len() > 0 {
		t.heap, v = t.heap[:t.Len()-1], t.heap[t.Len()-1]
	}
	r, ok := v.(Task)
	if !ok && v != nil {
		panic("node must type of Task")
	}
	return r
}

func (t *ComparablePipe) Less(i, j int) bool {
	return t.heap[i].CompareTo(t.heap[j].Val()) < 0
}

func (t *ComparablePipe) Swap(i, j int) {
	t.heap[i], t.heap[j] = t.heap[j], t.heap[i]
}

func (t *ComparablePipe) Len() int {
	return len(t.heap)
}

type DefaultTask struct {
	ctx   context.Context
	count int
}

func (t *DefaultTask) Run() {
	if t.count == 1 {
		time.Sleep(time.Duration(10) * time.Second)
	} else {
		time.Sleep(time.Duration(3) * time.Second)
	}
	println("===>Run：", t.count)
}
func (t *DefaultTask) Ctx() context.Context {
	return t.ctx
}

func (t *DefaultTask) Val() interface{} {
	return t.count
}
func (t *DefaultTask) CompareTo(other interface{}) int {
	i := other.(int) - t.count
	return i
}

func TestMultiCase(t *testing.T) {
	c0 := make(chan string)
	c1 := make(chan string)
	c2 := make(chan string)
	c3 := make(chan string)
	go func() {
		v := <-c3
		println(v)
	}()
	go func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		utils.NewMulti[string]("demo", ctx, logrus.New()).
			ChannelHandler(c0, func(val string) {
				println("channel0", val)
			}).ChannelHandler(c1, func(val string) {
			println("channel1", val)
		}).ChannelHandler(c2, func(val string) {
			println("channel2", val)
		}).ChannelSend(c3, func() string {
			return "sendVal"
		}).Default(func(t string) {
			println("===>default")
		}).Start()
	}()
	c0 <- "name1"
	go func() {
		c1 <- "name2"
	}()
	go func() {
		v := <-c3
		println(v)
	}()
	time.Sleep(20 * time.Second)
}

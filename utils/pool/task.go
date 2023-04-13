package pool

import (
	"context"
	"sync/atomic"
)

type Task interface {
	Run()
	Ctx() context.Context
	Interrupt() context.CancelFunc
}

type TaskPipe interface {
	PushTask(task Task)
	PopTask() Task
}

func GetFifoPipe(cap uint32) *FifoPipe {
	r := &FifoPipe{Cap: cap, Channel: make(chan Task, cap)}
	return r
}

type FifoPipe struct {
	Cap     uint32
	Channel chan Task
}

type FifoLine struct {
	task Task
	next *FifoLine
}

func (t *FifoPipe) PushTask(task Task) {
	t.CompareCall(1, func() {
		t.Channel <- task
	})
}

func (t *FifoPipe) PopTask() Task {
	if t.CompareCall(-1, func() {}) > 0 {
		return <-t.Channel
	}
	return nil
}

func (t *FifoPipe) CompareCall(payload int32, call func()) int {
	cap32 := int32(t.Cap)
	for {
		newVal := cap32 + payload
		if atomic.CompareAndSwapInt32(&cap32, cap32, newVal) {
			call()
			return int(newVal)
		}
	}
}

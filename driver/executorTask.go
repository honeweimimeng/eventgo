package driver

import (
	"context"
	"github.com/sirupsen/logrus"
	"sync"
)

type ExecutorTask struct {
	executor Executor
	channel  Channel
}

var (
	pool1 = sync.Pool{New: func() any {
		return &ExecutorTask{}
	}}
)

func NewTask(ex Executor, ch Channel) *ExecutorTask {
	task := pool1.Get().(*ExecutorTask)
	task.executor = ex
	task.channel = ch
	return task
}

func (e *ExecutorTask) Run() {
	e.executor.Execute(e.channel)
	pool1.Put(e)
}

func (e *ExecutorTask) Ctx() context.Context {
	return e.executor.Context().Context()
}

func (e *ExecutorTask) Interrupt() context.CancelFunc {
	e.Logger().Println("event executor task Interrupt Name:", e.executor.Name())
	return e.executor.Context().Interrupt()
}

func (e *ExecutorTask) Logger() logrus.StdLogger {
	return e.executor.Context().Config().Logger
}

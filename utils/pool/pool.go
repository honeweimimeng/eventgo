package pool

import (
	"context"
	"sync/atomic"
)

type Pool interface {
	StartUp()
	Name() string
	Run(task Task)
	WorkerCount() uint32
	Cap() uint32
	ExHandler(task Task, err error)
}

type DefaultPool struct {
	Name_     string
	workCount uint32
	Cap_      uint32
	Ctx       context.Context
	Pipe      TaskPipe
}

func (p *DefaultPool) StartUp() {
	go func() {
		for {
			select {
			case <-p.Ctx.Done():
				println("done msg:", p.Ctx.Err().Error())
				return
			default:
				atomic.AddUint32(&p.workCount, 1)
				if atomic.LoadUint32(&p.workCount) <= atomic.LoadUint32(&p.Cap_) {
					w := &Worker{Next: func() Task {
						return p.Pipe.PopTask()
					}, pool: p}
					w.doWork(p.Pipe.PopTask())
				}
			}
		}
	}()
}

func (p *DefaultPool) Name() string {
	return p.Name_
}

func (p *DefaultPool) Run(task Task) {
	p.Pipe.PushTask(task)
}

func (p *DefaultPool) WorkerCount() uint32 {
	return p.workCount
}

func (p *DefaultPool) Cap() uint32 {
	return p.Cap_
}

func (p *DefaultPool) ExHandler(task Task, err error) {
	println(task, "is failed,message: ", err.Error())
	if task != nil {
		cancel := task.Interrupt()
		cancel()
	}
}

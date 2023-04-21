package utils

import (
	"context"
	"github.com/honeweimimeng/eventgo/utils/pool"
)

type FixPool struct {
	ctx    context.Context
	cancel context.CancelFunc
	pool   pool.Pool
}

func NewFixPool(ctx context.Context, cancel context.CancelFunc, cap uint32, name string) *FixPool {
	core := &pool.DefaultPool{
		Cap_:  cap,
		Name_: name + "fixPool",
		Ctx:   ctx,
		Pipe:  pool.GetFifoPipe(cap),
	}
	return &FixPool{ctx: ctx, cancel: cancel, pool: core}
}

func (p *FixPool) GetFixPool() chan func() {
	runTask := make(chan func())
	p.pool.Run(NewFixPoolTask(p.ctx, p.cancel, runTask))
	return runTask
}

func (p *FixPool) Start() {
	p.pool.StartUp()
}

type FixPoolTask struct {
	cancelFunc context.CancelFunc
	ctx        context.Context
	chain      chan func()
}

func NewFixPoolTask(ctx context.Context, cancel context.CancelFunc, chain chan func()) *FixPoolTask {
	return &FixPoolTask{ctx: ctx, cancelFunc: cancel, chain: chain}
}

func (e *FixPoolTask) Run() {
	for {
		select {
		case f := <-e.chain:
			f()
		}
	}
}

func (e *FixPoolTask) Ctx() context.Context {
	return e.ctx
}

func (e *FixPoolTask) Interrupt() context.CancelFunc {
	return e.cancelFunc
}

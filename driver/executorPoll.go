package driver

import (
	"github.com/honeweimimeng/eventgo/utils"
	"github.com/honeweimimeng/eventgo/utils/pool"
)

type ExecutorPool struct {
	taskJoin        chan *ExecutorTask
	ctx             ExecutorContext
	pool            pool.Pool
	executorFutures []*utils.Future
}

func NewExecutorPoll(c ExecutorContext) *ExecutorPool {
	res := &ExecutorPool{
		ctx:      c,
		taskJoin: make(chan *ExecutorTask, c.Config().ExecutorCap),
		pool: &pool.DefaultPool{
			Cap_:  c.Config().ExecutorCap,
			Name_: c.Config().Name,
			Ctx:   c.Context(),
			Pipe:  pool.GetFifoPipe(c.Config().ExecutorCap),
		},
	}
	res.pool.StartUp()
	return res
}

func (p *ExecutorPool) Channel(executor Executor) chan *ExecutorTask {
	return nil
}

func (p *ExecutorPool) Join(task *ExecutorTask) ExecutorGroup {
	p.executorFutures = append(p.executorFutures, task.channel.Future())
	p.taskJoin <- task
	return p
}

func (p *ExecutorPool) Execute() []*utils.Future {
	go func() {
		for {
			select {
			case task := <-p.taskJoin:
				p.pool.Run(task)
			}
		}
	}()
	return p.executorFutures
}

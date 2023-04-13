package driver

import (
	"context"
	"github.com/honeweimimeng/eventgo/utils"
	"github.com/sirupsen/logrus"
)

type Channel interface {
	Future() *utils.Future
	Msg() any
}

type Executor interface {
	Name() string
	Execute(ch Channel)
	Context() ExecutorContext
}

type ExecutorContext interface {
	Config() *ExecutorConfig
	Group() ExecutorGroup
	SetGroup(group ExecutorGroup)
	Context() context.Context
	Interrupt() context.CancelFunc
	Process() ExecutorProcess
}

type ExecutorGroup interface {
	Join(task *ExecutorTask) ExecutorGroup
	Execute() []*utils.Future
	Channel(executor Executor) chan *ExecutorTask
}

type ExecutorProcess interface {
	Context(ctx ExecutorContext)
	Process(executor Executor)
	Group(group ExecutorGroup)
}

type ExecutorConfig struct {
	ExecutorCap uint32
	Name        string
	Logger      logrus.StdLogger
}

func DefaultConfig() *ExecutorConfig {
	return &ExecutorConfig{
		ExecutorCap: 1000,
		Name:        "exConfig0",
		Logger:      logrus.New(),
	}
}

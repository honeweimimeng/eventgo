package event

import (
	"github.com/honeweimimeng/eventgo/driver"
	"github.com/honeweimimeng/eventgo/utils"
)

type Proto interface {
	Name() string
	Id() int
	Channel() driver.Channel
}

type Handler interface {
	driver.Executor
	Events() []Proto
}

type LoopExecutor struct {
	name           string
	handles        []Handler
	context        driver.ExecutorContext
	sel            *utils.MultiCaseSel[*driver.ExecutorTask]
	eventHandleMap *utils.SafeMap[string, Handler]
}

func NewEventLoop(ctx driver.ExecutorContext, handles []Handler) *LoopExecutor {
	return (&LoopExecutor{
		name:           FormatHandleName(handles),
		handles:        handles,
		context:        ctx,
		sel:            utils.NewMulti[*driver.ExecutorTask](ctx.Config().Name, ctx.Context(), ctx.Config().Logger),
		eventHandleMap: utils.NewSafeMap[string, Handler](),
	}).initEventHandle()
}

func (l *LoopExecutor) initEventHandle() *LoopExecutor {
	for _, handle := range l.handles {
		for _, event := range handle.Events() {
			l.eventHandleMap.Put(event.Name(), handle)
		}
	}
	return l
}

func (l *LoopExecutor) EventHandle() *utils.SafeMap[string, Handler] {
	return l.eventHandleMap
}

func (l *LoopExecutor) AddTrigger(trigger Trigger) *LoopExecutor {
	bus := l.context.Group().(*Bus)
	bus.AddTrigger(trigger)
	return l
}

func (l *LoopExecutor) Task(ch driver.Channel) *driver.ExecutorTask {
	return driver.NewTask(l, ch)
}

func (l *LoopExecutor) TaskMsg(msg any) *driver.ExecutorTask {
	c := &Channel{MsgInstance: msg}
	return driver.NewTask(l, c)
}

func (l *LoopExecutor) Name() string {
	return l.name
}

func (l *LoopExecutor) Execute(ch driver.Channel) {
	for _, item := range l.handles {
		l.sel.ChannelHandler(l.context.Group().Channel(item), func(ex *driver.ExecutorTask) {
			l.Context().Group().Join(ex)
		})
	}
	l.context.Config().Logger.Println(l.Name())
	ch.Future().Complete()
	l.sel.Start()
}

func (l *LoopExecutor) Context() driver.ExecutorContext {
	return l.context
}

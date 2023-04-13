package event

import (
	"github.com/honeweimimeng/eventgo/driver"
	"github.com/honeweimimeng/eventgo/utils"
	"reflect"
)

type Bus struct {
	group            driver.ExecutorGroup
	trigger          Trigger
	adviceCh         chan []Proto
	ctx              driver.ExecutorContext
	registeredHandle []driver.Executor
	handleChannel    *utils.SafeMap[Handler, chan *driver.ExecutorTask]
	eventHandle      *utils.SafeMap[reflect.Type, Handler]
	camp             utils.CompositeFuture
}

func UseEventBus(ctx driver.ExecutorContext) *Bus {
	res := &Bus{
		registeredHandle: make([]driver.Executor, 0),
		group:            driver.NewExecutorPoll(ctx),
		trigger:          NewTriggerManager(ctx),
		adviceCh:         make(chan []Proto),
		ctx:              ctx,
		handleChannel:    utils.NewSafeMap[Handler, chan *driver.ExecutorTask](),
		eventHandle:      utils.NewSafeMap[reflect.Type, Handler](),
	}
	ctx.SetGroup(res)
	p := ctx.Process()
	if p != nil {
		p.Context(ctx)
	}
	return res
}

func (b *Bus) Join(task *driver.ExecutorTask) driver.ExecutorGroup {
	return b.group.Join(task)
}

func (b *Bus) Execute() []*utils.Future {
	return b.startEventLoop(b.group.Execute())
}

func (b *Bus) startEventLoop(futures []*utils.Future) []*utils.Future {
	started := &utils.Future{}
	res := []*utils.Future{started}
	(&b.camp).All(futures...).Result(func(r *utils.Result) {
		if !r.IsSuccess() {
			b.ctx.Config().Logger.Panicln("cannot start eventLoop,because executor has failed")
		}
		b.trigger.AcceptEvents(b.adviceCh)
		go func() {
			for {
				select {
				case acceptedEvents := <-b.adviceCh:
					b.processEvents(acceptedEvents)
					break
				}
			}
		}()
	})
	return res
}

func (b *Bus) registerEventHandler(handler Handler) *Bus {
	for _, item := range handler.Events() {
		b.eventHandle.Put(reflect.TypeOf(item), handler)
	}
	b.registeredHandle = append(b.registeredHandle, handler)
	return b
}

func (b *Bus) Channel(executor driver.Executor) chan *driver.ExecutorTask {
	handler, ok := executor.(Handler)
	if ok {
		ch := b.handleChannel.Get(handler)
		if ch == nil {
			ch = make(chan *driver.ExecutorTask)
			b.handleChannel.Put(handler, ch)
		}
		b.registerEventHandler(handler)
		return ch
	}
	b.ctx.Config().Logger.Panicln("cannot init channel,because executor not type of handler")
	return nil
}

func (b *Bus) AddTrigger(trigger Trigger) *Bus {
	b.trigger.Child(trigger)
	return b
}

func (b *Bus) processEvents(events []Proto) {
	for _, item := range events {
		ex := b.eventHandle.Get(reflect.TypeOf(item))
		b.ctx.Group().Channel(ex) <- driver.NewTask(ex, item.Channel())
	}
}

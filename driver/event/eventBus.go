package event

import (
	"github.com/honeweimimeng/eventgo/driver"
	"github.com/honeweimimeng/eventgo/utils"
)

type Bus struct {
	eventProcessPool *utils.FixPool
	group            driver.ExecutorGroup
	trigger          Trigger
	adviceCh         chan []Proto
	ctx              driver.ExecutorContext
	registeredHandle []driver.Executor
	handleChannel    *utils.SafeMap[Handler, chan *driver.ExecutorTask]
	eventHandle      *utils.SafeMap[string, Handler]
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
		eventHandle:      utils.NewSafeMap[string, Handler](),
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
		var chains []TriggerAdvice
		if manager, ok := b.trigger.(*TriggerManager); ok {
			chains = manager.AcceptMultiTrigger()
		} else {
			b.trigger.AcceptEvents(b.adviceCh)
			chains = []TriggerAdvice{{protoChain: b.adviceCh, global: b.trigger.Global()}}
		}
		b.InitProcessEventLen(len(chains)).Run(chains)
	})
	return res
}

func (b *Bus) InitProcessEventLen(cap int) *Bus {
	ctx := b.ctx
	b.eventProcessPool = utils.NewFixPool(ctx.Context(), ctx.Interrupt(), uint32(cap), ctx.Config().Name)
	return b
}

func (b *Bus) Run(chains []TriggerAdvice) {
	sel := utils.NewMulti[[]Proto]("eventBus sel-case loop", b.ctx.Context(), b.ctx.Config().Logger)
	for _, item := range chains {
		fixProcess := b.eventProcessPool.GetFixPool()
		handleMap := b.eventHandle
		if !item.global && item.loopGroup != nil {
			handleMap = item.loopGroup.EventHandle()
		}
		sel.ChannelHandler(item.protoChain, func(acceptedEvents []Proto) {
			fixProcess <- func() { b.processEvents(acceptedEvents, handleMap) }
		})
	}
	b.eventProcessPool.Start()
	sel.StartAsync()
}

func (b *Bus) registerEventHandler(handler Handler) *Bus {
	for _, item := range handler.Events() {
		b.eventHandle.Put(item.Name(), handler)
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

func (b *Bus) processEvents(events []Proto, eventHandle *utils.SafeMap[string, Handler]) {
	for _, item := range events {
		ex := eventHandle.Get(item.Name())
		if ex == nil {
			b.ctx.Config().Logger.Println("cannot process event", item.Name(),
				", because that is not registered in handle and trigger is not global trigger")
			continue
		}
		b.ctx.Group().Channel(ex) <- driver.NewTask(ex, item.Channel())
	}
}

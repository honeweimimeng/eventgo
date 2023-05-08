package event

import (
	"github.com/honeweimimeng/eventgo/driver"
	"github.com/honeweimimeng/eventgo/utils"
)

type BootStrap struct {
	loop []*LoopDesc
}

func (b *BootStrap) EventLoop() *LoopDesc {
	return &LoopDesc{strap: b}
}

func (b *BootStrap) Context(ctx driver.ExecutorContext) {
	for _, item := range b.loop {
		item.SetCtx(ctx)
	}
}

func (b *BootStrap) Process(executor driver.Executor) {}

func (b *BootStrap) Group(group driver.ExecutorGroup) {}

func (b *BootStrap) StartUp() []*utils.Future {
	ctx := Process(b).LoadProperty(driver.DefaultConfig())
	return UseEventBus(ctx).Execute()
}

type LoopDesc struct {
	trigger *ProxyTrigger
	handles []*ProxyHandler
	strap   *BootStrap
}

func (l *LoopDesc) Handle(events []Proto, handle SimpleHandler) *LoopDesc {
	proxy := &ProxyHandler{events: events, handle: handle}
	l.handles = append(l.handles, proxy)
	return l
}

func (l *LoopDesc) Handler(name string, events []Proto, handle SimpleHandler) *LoopDesc {
	proxy := &ProxyHandler{name: name, events: events, handle: handle}
	l.handles = append(l.handles, proxy)
	return l
}

func (l *LoopDesc) Trigger(trigger SimpleTrigger, isInitializer bool) *BootStrap {
	l.trigger = &ProxyTrigger{trigger: trigger, global: true, isInitializer: isInitializer}
	return l.boot()
}

func (l *LoopDesc) ExTrigger(trigger SimpleTrigger, isInitializer bool) *BootStrap {
	l.trigger = &ProxyTrigger{trigger: trigger, global: false, isInitializer: isInitializer}
	return l.boot()
}

func (l *LoopDesc) MultiTrigger(trigger ...SimpleTrigger) *LoopDesc {
	var head *ProxyTrigger
	cur := head
	for _, item := range trigger {
		if head == nil {
			head = &ProxyTrigger{trigger: item}
			continue
		}
		cur.next = &ProxyTrigger{trigger: item}
		cur = cur.next.(*ProxyTrigger)
	}
	l.trigger = head
	return l
}

func (l *LoopDesc) SetCtx(ctx driver.ExecutorContext) {
	var handleInters []Handler
	for _, h := range l.handles {
		h.ctx = ctx
		handleInters = append(handleInters, h)
	}
	eventLoop := NewEventLoop(ctx, handleInters)
	l.trigger.ctx = ctx
	l.trigger.group = eventLoop
	task := eventLoop.AddTrigger(l.trigger).TaskMsg("event loop")
	ctx.Group().Join(task)
}

func (l *LoopDesc) boot() *BootStrap {
	l.strap.loop = append(l.strap.loop, l)
	return l.strap
}

type SimpleHandler func(ch driver.Channel)

type SimpleTrigger func(ch chan []Proto)

type ProxyHandler struct {
	ctx    driver.ExecutorContext
	name   string
	events []Proto
	handle SimpleHandler
}

func (h *ProxyHandler) Execute(ch driver.Channel) {
	h.handle(ch)
}

func (h *ProxyHandler) Events() []Proto {
	return h.events
}

func (h *ProxyHandler) Name() string {
	return h.name
}

func (h *ProxyHandler) Context() driver.ExecutorContext {
	return h.ctx
}

type ProxyTrigger struct {
	global        bool
	ctx           driver.ExecutorContext
	next          Trigger
	trigger       SimpleTrigger
	group         *LoopExecutor
	isInitializer bool
}

func (e *ProxyTrigger) AcceptEvents(ch chan []Proto) {
	e.trigger(ch)
}

func (e *ProxyTrigger) Next() Trigger {
	return e.next
}

func (e *ProxyTrigger) Child(trigger Trigger) {
	e.next = trigger
}

func (e *ProxyTrigger) Global() bool {
	return e.global
}

func (e *ProxyTrigger) LoopGroup() *LoopExecutor {
	return e.group
}

func (e *ProxyTrigger) Initializer() bool {
	return e.isInitializer
}

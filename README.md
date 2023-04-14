#### eventGo
![](https://raw.githubusercontent.com/honeweimimeng/eventgo/master/icons/progress-100-green.jpg?raw=true)
![](https://github.com/honeweimimeng/eventgo/blob/master/icons/gopool-100%25-green.jpg?raw=true)
![](https://raw.githubusercontent.com/honeweimimeng/eventgo/master/icons/Variableselect--case-100%25-green.jpg)
- High performance event driver framework ⌛
- Flexible Event Notification 📩
- Controllable Event Channel 📦
- configurable Goroutine pool 🏂
### Demo
```golang
var boot event.BootStrap
boot.EventLoop().Handle(Events(0), func(ch driver.Channel) {
    println("===>")
	time.Sleep(2 * time.Second)
	println("===>finish")
}).Trigger(func(ch chan []event.Proto) {
	ch <- []event.Proto{HappenEvent(0, "WRITE", "write msg")}
	ch <- []event.Proto{HappenEvent(0, "WRITE", "write msg")}
}).Boot().StartUp()
time.Sleep(200 * time.Second)
```
### BootStrap
> BootStrap provides a fast and elegant way to create,zero friendly
### Advanced Usage
#### Step1 Create Event Handle
> Inheriting event.Handler is to implement Execute
```golang
type ProxyHandler struct {
    ctx    driver.ExecutorContext
    name   string
    events []Proto
    handle func....
}

func (h *ProxyHandler) Execute(ch driver.Channel) {
    func.....
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
```
#### Step2 Create EventTrigger
> Inherit EventTrigger to implement event notifier
```golang
type ProxyTrigger struct {
    ctx     driver.ExecutorContext
    next    Trigger
    trigger SimpleTrigger
}

func (e *ProxyTrigger) AcceptEvents(ch chan []Proto) {
	//Use channel to send events
    e.trigger(ch)
}

func (e *ProxyTrigger) Next() Trigger {
    return e.next
}

func (e *ProxyTrigger) Child(trigger Trigger) {
    e.next = trigger
}
```
#### Step3 Register them
> Register to event group by ExecutorContext in your ExecutorProcess
```golang
type BootStrap struct {
    loop []*LoopDesc
}

func (b *BootStrap) EventLoop() *LoopDesc {
    return &LoopDesc{strap: b}
}

func (b *BootStrap) Context(ctx driver.ExecutorContext) {
    ctx.Group().Join(.....)
}
```
#### Step4 use them
```golang
ctx := event.Process(&BootStrap{}).LoadProperty(driver.DefaultConfig())
event.UseEventBus(ctx).Execute()
time.Sleep(200 * time.Second)
```
>Tips:BootStrap is actually a special form of advanced usage, you can refer to Implement your own event-driven
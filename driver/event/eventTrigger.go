package event

import (
	"github.com/honeweimimeng/eventgo/driver"
	"github.com/honeweimimeng/eventgo/utils"
)

type TriggerAdvice struct {
	protoChain chan []Proto
	global     bool
	loopGroup  *LoopExecutor
}

type Trigger interface {
	Global() bool
	AcceptEvents(ch chan []Proto)
	Next() Trigger
	Child(trigger Trigger)
	LoopGroup() *LoopExecutor
}

type TriggerManager struct {
	ctx      driver.ExecutorContext
	childCur int
	child    []Trigger
	sel      *utils.MultiCaseSel[[]Proto]
}

func NewTriggerManager(ctx driver.ExecutorContext) *TriggerManager {
	res := &TriggerManager{
		sel: utils.NewMulti[[]Proto](ctx.Config().Name, ctx.Context(), ctx.Config().Logger),
		ctx: ctx,
	}
	return res
}

func (m *TriggerManager) listenEvent(ch chan []Proto) (*TriggerManager, []TriggerAdvice) {
	var chains []TriggerAdvice
	for item := m.Next(); item != nil; item = m.Next() {
		chItem := ch
		if chItem == nil {
			chItem = make(chan []Proto)
			triggerAdvice := TriggerAdvice{protoChain: chItem, global: item.Global(), loopGroup: item.LoopGroup()}
			chains = append(chains, triggerAdvice)
		}
		m.processChildNext(chItem, item)
	}
	m.sel.StartAsync()
	return m, chains
}

func (m *TriggerManager) processChildNext(ch chan []Proto, trigger Trigger) {
	go trigger.AcceptEvents(ch)
	for item := trigger.Next(); item != nil; item = item.Next() {
		item.AcceptEvents(ch)
	}
}

func (m *TriggerManager) AcceptMultiTrigger() []TriggerAdvice {
	_, chains := m.listenEvent(nil)
	return chains
}

func (m *TriggerManager) Global() bool {
	return true
}

func (m *TriggerManager) AcceptEvents(ch chan []Proto) {
	m.listenEvent(ch)
}

func (m *TriggerManager) Next() Trigger {
	if len(m.child) < m.childCur+1 {
		return nil
	}
	trigger := m.child[m.childCur]
	m.childCur = m.childCur + 1
	return trigger
}

func (m *TriggerManager) Child(trigger Trigger) {
	m.child = append(m.child, trigger)
}

func (m *TriggerManager) LoopGroup() *LoopExecutor {
	return nil
}

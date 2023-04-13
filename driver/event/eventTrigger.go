package event

import (
	"github.com/honeweimimeng/eventgo/driver"
	"github.com/honeweimimeng/eventgo/utils"
)

type Trigger interface {
	AcceptEvents(ch chan []Proto)
	Next() Trigger
	Child(trigger Trigger)
}

type TriggerManager struct {
	ctx  driver.ExecutorContext
	next Trigger
	sel  *utils.MultiCaseSel[[]Proto]
	hand func(proto []Proto)
}

func NewTriggerManager(ctx driver.ExecutorContext) *TriggerManager {
	res := &TriggerManager{
		sel: utils.NewMulti[[]Proto](ctx.Config().Name, ctx.Context(), ctx.Config().Logger),
		ctx: ctx,
	}
	return res
}

func (m *TriggerManager) listenEvent() *TriggerManager {
	go func() {
		for i := m.Next(); i != nil; i = i.Next() {
			chItem := make(chan []Proto)
			go i.AcceptEvents(chItem)
			m.sel.ChannelHandler(chItem, func(proto []Proto) {
				m.hand(proto)
			})
		}
		m.sel.Start()
	}()
	return m
}

func (m *TriggerManager) AcceptEvents(ch chan []Proto) {
	m.listenEvent()
	if m.hand == nil {
		m.hand = func(proto []Proto) {
			ch <- proto
		}
	}
}

func (m *TriggerManager) Next() Trigger {
	return m.next
}

func (m *TriggerManager) Child(trigger Trigger) {
	trigger.Child(m.next)
	m.next = trigger
}

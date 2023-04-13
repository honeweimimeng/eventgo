package event

import "github.com/honeweimimeng/eventgo/utils"

type Channel struct {
	EventProto  Proto
	MsgInstance any
	Future_     utils.Future
}

func (c *Channel) Msg() any {
	return c.MsgInstance
}

func (c *Channel) Future() *utils.Future {
	return &c.Future_
}

func (c *Channel) Proto() Proto {
	return c.EventProto
}

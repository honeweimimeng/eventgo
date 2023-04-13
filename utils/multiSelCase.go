package utils

import (
	"context"
	"github.com/sirupsen/logrus"
	"reflect"
)

type MultiCaseSel[T any] struct {
	Name          string
	Ctx           context.Context
	cases         []reflect.SelectCase
	Logger        logrus.StdLogger
	handleMapping map[reflect.SelectCase]func(T)
}

func NewMulti[T any](name string, ctx context.Context, logger logrus.StdLogger) *MultiCaseSel[T] {
	res := &MultiCaseSel[T]{
		Name:   name,
		Ctx:    ctx,
		Logger: logger,
	}
	cases := []reflect.SelectCase{{
		Dir:  reflect.SelectRecv,
		Chan: reflect.ValueOf(res.Ctx.Done()),
	}}
	res.cases = cases
	res.handleMapping = make(map[reflect.SelectCase]func(T))
	return res
}

func (sel *MultiCaseSel[T]) ChannelHandler(channel chan T, handle func(T)) *MultiCaseSel[T] {
	newCase := reflect.SelectCase{
		Dir:  reflect.SelectRecv,
		Chan: reflect.ValueOf(channel),
	}
	sel.cases = append(sel.cases, newCase)
	sel.handleMapping[newCase] = handle
	return sel
}

func (sel *MultiCaseSel[T]) ChannelSend(channel chan T, handle func() T) *MultiCaseSel[T] {
	newCase := reflect.SelectCase{
		Dir:  reflect.SelectSend,
		Chan: reflect.ValueOf(channel),
		Send: reflect.ValueOf(handle()),
	}
	sel.cases = append(sel.cases, newCase)
	sel.handleMapping[newCase] = func(t T) {}
	return sel
}

func (sel *MultiCaseSel[T]) ChannelSendProcess(channel chan T, process func(t T), handle func() T) *MultiCaseSel[T] {
	newCase := reflect.SelectCase{
		Dir:  reflect.SelectSend,
		Chan: reflect.ValueOf(channel),
		Send: reflect.ValueOf(handle()),
	}
	sel.cases = append(sel.cases, newCase)
	sel.handleMapping[newCase] = process
	return sel
}

func (sel *MultiCaseSel[T]) Default(handle func(T)) *MultiCaseSel[T] {
	newCase := reflect.SelectCase{
		Dir: reflect.SelectDefault,
	}
	sel.cases = append(sel.cases, newCase)
	sel.handleMapping[newCase] = handle
	return sel
}

func (sel *MultiCaseSel[T]) Start() {
	for {
		chosen, val, ok := reflect.Select(sel.cases)
		dir := sel.cases[chosen].Dir
		if o := dir == reflect.SelectSend || dir == reflect.SelectDefault; ok || o {
			v, bre := sel.covetT(ok, val)
			if bre {
				return
			}
			sel.handleMapping[sel.cases[chosen]](v)
			continue
		}
		sel.Logger.Fatal("select ", sel.Name, " context cancel")
		return
	}
}

func (sel *MultiCaseSel[T]) covetT(isRec bool, val reflect.Value) (T, bool) {
	defer func() {
		if r := recover(); r != nil && isRec {
			sel.Logger.Fatal("unexpect data type in select ", sel.Name)
		}
	}()
	intVal, okVal := val.Interface().(T)
	return intVal, !okVal
}

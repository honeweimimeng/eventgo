package utils

import "testing"

func TestAdvice(t *testing.T) {
	f := &Future{}
	go func() {
		f.Complete()
	}()
	f.Success(func() {

	}).Fail(func() {

	})
	f1 := &Future{}
	go func() {
		f1.Complete()
	}()
	f1.Success(func() {

	}).Fail(func() {

	})
	var c CompositeFuture
	c.All(f, f1)
	c.Result(func(r *Result) {
		println(r.IsSuccess())
	})
}

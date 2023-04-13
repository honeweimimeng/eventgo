package utils

import "sync"

type Call func()

type Future struct {
	success Call
	failed  Call
	isWait  bool
	wait    sync.WaitGroup
}

func (f *Future) Get() {
	if f.isWait {
		f.wait.Done()
	}
}

func (f *Future) Fail(call Call) *Future {
	f.failed = call
	return f
}

func (f *Future) Success(call Call) *Future {
	f.success = call
	return f
}

func (f *Future) finish(wait bool) {
	if wait && !f.isWait {
		f.isWait = true
		f.wait.Add(1)
		f.wait.Wait()
	}
}

func (f *Future) Complete() {
	f.finish(f.success == nil)
	f.success()
}

func (f *Future) Failed() {
	f.finish(f.failed == nil)
	f.failed()
}

type Result struct {
	isFailed bool
}

func (r *Result) IsSuccess() bool {
	return !r.isFailed
}

type CompositeFuture struct {
	r    Result
	wait sync.WaitGroup
}

func (c *CompositeFuture) All(f ...*Future) *CompositeFuture {
	for _, item := range f {
		item.Success(func() {
			c.wait.Done()
		}).Fail(func() {
			c.wait.Done()
			c.r.isFailed = true
		}).Get()
	}
	c.wait.Add(len(f))
	return c
}

func (c *CompositeFuture) Result(f func(r *Result)) {
	c.wait.Wait()
	f(&c.r)
}

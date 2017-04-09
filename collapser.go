// Copyright 2017 Cristian Greco <cristian@regolo.cc>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package collapser provides a function call deduplication utility.
//
// In many real-world scenarios, under high load or traffic conditions,
// an application may end up with several goroutines making the very same
// - expensive - function call at the same time. A Collapser provides a lock
// mechanism that enables multiple identical requests to be processed as one
// single request: while the first goroutine is making progress, the rest will
// be blocked waiting for the collapsed function to complete, and then will
// receive the same result.
package collapser

import (
	"sync"
	"sync/atomic"
)

// A Collapser prevents duplicate function calls to run simultaneously.
// If any identical function invocations occur while the original function
// is still executing, the callers will wait for the original function to
// complete execution and then the result will be shared between callers.
type Collapser struct {
	m     sync.Mutex
	tasks map[interface{}]*task
}

// TaskResult provides access to the function result.
type TaskResult struct {
	res interface{}
	cnt uint64
}

// Get returns the result of the original function invocation.
func (r *TaskResult) Get() interface{} {
	return r.res
}

// Collapsed returns the number of invocations that have been collapsed.
func (r *TaskResult) Collapsed() uint64 {
	return atomic.LoadUint64(&r.cnt)
}

// task represents a function invocation, it implements a mechanism to
// block simultaneous invocations and provides access to shared result.
type task struct {
	rw  sync.RWMutex
	res *TaskResult
}

// NewCollapser returns a new Collapser in initialized state.
func NewCollapser() *Collapser {
	return &Collapser{
		tasks: make(map[interface{}]*task),
	}
}

// Do executes the given function and returns a TaskResult object wrapping
// the function result. The given key is used by the Collapser to uniquely
// identify the function invocation. Any identical function invocation (i.e.
// identified by the same key) which occurs while the original function is
// executing will be blocked waiting for the running invocation to complete,
// and will receive a TaskResult object wrapping the function result of the
// original function.
func (c *Collapser) Do(key interface{}, f func() interface{}) *TaskResult {
	c.m.Lock()
	t, exists := c.tasks[key]
	if !exists {
		t = &task{}
		t.res = &TaskResult{}
		t.rw.Lock()
		c.tasks[key] = t
	}
	c.m.Unlock()

	atomic.AddUint64(&t.res.cnt, 1)

	if exists {
		t.rw.RLock()
		t.rw.RUnlock()
	} else {
		t.res.res = f()
		t.rw.Unlock()

		c.m.Lock()
		delete(c.tasks, key)
		c.m.Unlock()
	}

	return t.res
}

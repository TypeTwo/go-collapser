// Copyright 2017 Cristian Greco <cristian@regolo.cc>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package collapser

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// Tests that Do() invokes the given function, and the returned TaskResult
// contains the function return value.
func TestDo(t *testing.T) {
	coll := NewCollapser()

	pre := 1
	val := 42

	res := coll.Do("key", func() interface{} {
		pre = val
		return val
	})

	if pre != val {
		t.Fatal("expected collapsed function to be invoked")
	}

	if res.Get().(int) != val {
		t.Fatal("expected collapser result to contain function result")
	}

	if res.Collapsed() != 1 {
		t.Fatal("expected collapsed count to be 1")
	}
}

// Tests that subsequent (serial) invocations of Do() execute func
// multiple times.
func TestDoInvokeSerial(t *testing.T) {
	c := NewCollapser()

	var cnt int

	f := func() interface{} {
		defer func() {
			cnt++
		}()
		return cnt
	}

	n := 3

	for i := 0; i < n; i++ {
		res := c.Do("key", f)
		if res.Get().(int) != i {
			t.Fatalf("expected invocation to return value %d", i)
		}
		if res.Collapsed() != 1 {
			t.Fatal("expected collapsed count to be 1")
		}
	}

	if cnt != n {
		t.Fatalf("expected collapsed function to be invoked %d times", n)
	}
}

// Tests that multiple (parallel) invocations of Do() execute func
// exactly once.
func TestDoInvokeParallel(t *testing.T) {
	c := NewCollapser()

	var text = "text"
	var cnt int

	f := func() interface{} {
		// sleep to give a chance to other
		// goroutines to enter c.Do()
		time.Sleep(100 * time.Millisecond)
		cnt++
		return text
	}

	n := 4
	wg := &sync.WaitGroup{}

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			res := c.Do("key", f)
			if res.Get().(string) != text {
				t.Fatal("expected invocation to return a string value")
			}
			if res.Collapsed() != uint64(n) {
				t.Fatalf("expected collapsed count to be %d", n)
			}
			wg.Done()
		}()
	}

	wg.Wait()

	if cnt != 1 {
		t.Fatal("expected collapsed function to be invoked once")
	}
}

// Tests that multiple (parallel) invocations of Do() using different keys
// are not collapsed.
func TestDoDistinctKeys(t *testing.T) {
	c := NewCollapser()

	var text = "text"
	var cnt uint64

	f := func() interface{} {
		// sleep to give a chance to other
		// goroutines to enter c.Do()
		time.Sleep(100 * time.Millisecond)
		atomic.AddUint64(&cnt, 1)
		return text
	}

	n := 4
	wg := &sync.WaitGroup{}

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			res := c.Do(i, f)
			if res.Get().(string) != text {
				t.Fatal("expected invocation to return a string value")
			}
			if res.Collapsed() != 1 {
				t.Fatal("expected collapsed count to be 1")
			}
			wg.Done()
		}(i)
	}

	wg.Wait()

	if got := int(atomic.LoadUint64(&cnt)); got != n {
		t.Fatalf("expected collapsed function to be invoked %d times", n)
	}
}

func benchmarkDo(n int, b *testing.B) {
	c := NewCollapser()

	var cnt uint64

	for i := 0; i < b.N; i++ {
		wg := &sync.WaitGroup{}
		for j := 0; j < n; j++ {
			wg.Add(1)
			go func() {
				c.Do("key", func() interface{} {
					return atomic.AddUint64(&cnt, 1)
				}).Get()
				wg.Done()
			}()
		}
		wg.Wait()
	}
}

func BenchmarkDo10(b *testing.B) {
	benchmarkDo(10, b)
}

func BenchmarkDo20(b *testing.B) {
	benchmarkDo(20, b)
}

package rps

import "sync/atomic"

type Counter struct {
	count atomic.Int64
}

func NewCounter() *Counter {
	return &Counter{
		count: atomic.Int64{},
	}
}

func (c *Counter) Inc() {
	c.count.Add(1)
}

func (c *Counter) Load() int64 {
	return c.count.Load()
}

func (c *Counter) Reset() {
	c.count.Store(0)
}

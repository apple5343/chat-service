package rps

import (
	"context"
	"time"
)

type RPSCounter struct {
	counter *Counter
	history []int64
}

func NewRPSCounter() *RPSCounter {
	return &RPSCounter{
		counter: NewCounter(),
		history: make([]int64, 0, 100),
	}
}

func (r *RPSCounter) Run(ctx context.Context) {
	ticker := time.NewTicker(time.Second)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				r.history = append(r.history, r.counter.Load())
				r.counter.Reset()
			}
		}
	}()
}

func (r *RPSCounter) Inc() {
	r.counter.Inc()
}

func (r *RPSCounter) History() []int64 {
	return r.history
}

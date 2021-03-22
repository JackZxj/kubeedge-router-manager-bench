package util

import "sync"

// Pool define the resource pool
type Pool struct {
	queue chan int
	wg    *sync.WaitGroup
}

// NewPool create a new resource pool
func NewPool(size int) *Pool {
	if size <= 0 {
		size = 1
	}
	return &Pool{
		queue: make(chan int, size),
		wg:    &sync.WaitGroup{},
	}
}

// Add acquire delta of channels
func (p *Pool) Add(delta int) {
	for i := 0; i < delta; i++ {
		p.queue <- 1
	}
	for i := 0; i > delta; i-- {
		<-p.queue
	}
	p.wg.Add(delta)
}

// Done release a channel
func (p *Pool) Done() {
	<-p.queue
	p.wg.Done()
}

// Wait for all channel
func (p *Pool) Wait() {
	p.wg.Wait()
}
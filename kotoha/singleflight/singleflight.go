package singleflight

import "sync"

type call[T any] struct {
	wg  sync.WaitGroup
	res T
	err error
}

type Flight[T any] struct {
	mu sync.Mutex
	m  map[string]*call[T]
}

// Do 泛型版的SingleFlight
func (g *Flight[T]) Do(key string, fn func() (T, error)) (T, error) {
	g.mu.Lock()
	if g.m == nil {
		g.m = make(map[string]*call[T])
	}
	if c, ok := g.m[key]; ok {
		g.mu.Unlock()
		c.wg.Wait()
		return c.res, c.err
	}
	c := new(call[T])
	c.wg.Add(1)
	g.m[key] = c
	g.mu.Unlock()

	c.res, c.err = fn()
	c.wg.Done()

	g.mu.Lock()
	delete(g.m, key)
	g.mu.Unlock()

	return c.res, c.err
}

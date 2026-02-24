package market

import (
	"context"
	"errors"
	"sync"
	"time"
)

type cache struct {
	mu       sync.Mutex
	cond     *sync.Cond
	ttl      time.Duration
	fetching bool
	expires  time.Time
	hasData  bool
	data     Snapshot
	lastErr  error
}

func newCache(ttl time.Duration) *cache {
	c := &cache{ttl: ttl}
	c.cond = sync.NewCond(&c.mu)
	return c
}

func (c *cache) get(ctx context.Context, fetch func(context.Context) (Snapshot, error)) (Snapshot, error) {
	for {
		c.mu.Lock()
		now := time.Now()
		if c.hasData && now.Before(c.expires) {
			data := c.data
			c.mu.Unlock()
			return data, nil
		}
		if c.fetching {
			c.cond.Wait()
			c.mu.Unlock()
			continue
		}
		c.fetching = true
		c.mu.Unlock()

		data, err := fetch(ctx)

		c.mu.Lock()
		c.fetching = false
		if err == nil {
			c.data = data
			c.hasData = true
			c.expires = time.Now().Add(c.ttl)
			c.lastErr = nil
		} else {
			c.lastErr = err
		}
		c.cond.Broadcast()
		c.mu.Unlock()

		if err == nil {
			return data, nil
		}
		if c.hasData {
			return c.data, nil
		}
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return Snapshot{}, err
		}
		return Snapshot{}, err
	}
}

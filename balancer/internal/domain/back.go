package domain

import (
	"net/url"
	"sync/atomic"
)

type Backend struct {
	URL        *url.URL
	ConnsCount int32
	Index      int
}

func (b *Backend) IncrementConns() {
	atomic.AddInt32(&b.ConnsCount, 1)
}

func (b *Backend) DecrementConns() {
	atomic.AddInt32(&b.ConnsCount, -1)
}

func (b *Backend) ResetConns() {
	atomic.StoreInt32(&b.ConnsCount, 0)
}

func (b *Backend) GetConns() int32 {
	return atomic.LoadInt32(&b.ConnsCount)
}

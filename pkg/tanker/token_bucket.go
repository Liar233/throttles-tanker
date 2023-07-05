package tanker

import (
	"sync/atomic"
)

type TokenBucket struct {
	capacity  uint32
	perUpdate uint32
	ticks     uint32
}

func NewTokenBucket(capacity, perUpdate uint32) *TokenBucket {

	return &TokenBucket{
		capacity:  capacity,
		perUpdate: perUpdate,
		ticks:     0,
	}
}

func (tb *TokenBucket) Fire() bool {

	if tb.ticks >= tb.capacity {

		return false
	}

	atomic.AddUint32(&tb.ticks, 1)

	return true
}

func (tb *TokenBucket) Update() uint32 {

	if atomic.LoadUint32(&tb.ticks) <= tb.perUpdate {

		atomic.StoreUint32(&tb.ticks, 0)

		return 0
	}

	atomic.StoreUint32(&tb.ticks, tb.ticks-tb.perUpdate)

	return tb.ticks
}

func (tb *TokenBucket) Rest() uint32 {

	return tb.capacity - atomic.LoadUint32(&tb.ticks)
}

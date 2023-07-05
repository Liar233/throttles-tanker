package tanker

import (
	"sync"
	"time"
)

type Tank struct {
	sync.Mutex
	buckets       map[string]Bucket
	update        time.Duration
	bucketBuilder func() Bucket
}

func (t *Tank) Check(id string) bool {

	t.Lock()

	defer t.Unlock()

	if bucket, ok := t.buckets[id]; ok {

		return bucket.Fire()
	}

	bucket := t.bucketBuilder()

	t.buckets[id] = bucket

	return t.buckets[id].Fire()
}

func (t *Tank) Flush() {

	t.Lock()

	defer t.Unlock()

	for key, bucket := range t.buckets {

		if bucket.Update() == 0 {

			delete(t.buckets, key)
		}
	}
}

func (t *Tank) Add(id string) Bucket {

	t.Lock()

	defer t.Unlock()

	bucket := t.bucketBuilder()

	t.buckets[id] = bucket

	return bucket
}

func (t *Tank) Remove(id string) {

	t.Lock()

	defer t.Unlock()

	delete(t.buckets, id)
}

func NewTank(builder func() Bucket) *Tank {

	return &Tank{
		buckets:       make(map[string]Bucket),
		bucketBuilder: builder,
	}
}

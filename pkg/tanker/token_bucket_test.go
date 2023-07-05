package tanker

import (
	"testing"
)

const tokenBucketCapacity = uint32(10)
const tokenBucketPerUpdate = uint32(2)

func TestNewTokenBucket(t *testing.T) {

	bucket := NewTokenBucket(10, 2)

	if bucket.ticks != 0 {

		t.Error("Created TokenBucket ticks not zero")
	}

	if bucket.perUpdate != tokenBucketPerUpdate {

		t.Error("Created TokenBucket perUpdate does not equal start value")
	}

	if bucket.capacity != tokenBucketCapacity {

		t.Error("Created TokenBucket capacity does not equal start value ")
	}
}

func TestTokenBucket_Fire(t *testing.T) {

	bucket := NewTokenBucket(tokenBucketCapacity, tokenBucketPerUpdate)

	for i := uint32(1); i < tokenBucketCapacity+1; i++ {

		if !bucket.Fire() {

			t.Error("TokenBucket Fire() method failed counting")
			t.Skip()

			break
		}

		if bucket.ticks != i {

			t.Error("TokenBucket Fire() not valid ticks increment")
			t.Skip()

			break
		}
	}

	if bucket.Fire() {

		t.Error("TokenBucket Fire() method failed counter overflow")
	}

}

func TestTokenBucket_Update(t *testing.T) {

	bucket := NewTokenBucket(tokenBucketCapacity, tokenBucketPerUpdate)

	bucket.Update()

	if bucket.ticks > 0 {

		t.Error("Created TokenBucket Update() method failed counter")
	}

	bucket.ticks = uint32(tokenBucketCapacity)

	bucket.Update()

	if bucket.ticks != tokenBucketCapacity-tokenBucketPerUpdate {

		t.Error("TokenBucket Update() method failed counter")
	}
}

func TestTokenBucket_Rest(t *testing.T) {

	bucket := NewTokenBucket(tokenBucketCapacity, tokenBucketPerUpdate)

	if bucket.Rest() != tokenBucketCapacity {

		t.Error("Created TokenBucket Rest() method not valid capacity value")
	}

	_ = bucket.Fire()

	if bucket.Rest() != tokenBucketCapacity-1 {

		t.Error("TokenBucket Rest() method not valid capacity value after Fire()")
	}
}

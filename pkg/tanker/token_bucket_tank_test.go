package tanker

import (
	"testing"
)

const capacitySize = uint32(2)
const perUpdateValue = uint32(1)

func TestTokenBucket_Tank_Compatibility(t *testing.T) {

	builder := func() Bucket {

		capacity := capacitySize
		perUpdate := perUpdateValue

		return NewTokenBucket(capacity, perUpdate)
	}

	tank := NewTank(builder)

	if !tank.Check("test") {

		t.Error("Tank Check() method not pass with TokenBucket implementation")
	}

	if !tank.Check("test") {

		t.Error("Tank Check() method not valid with TokenBucket implementation")
	}

	bucket, _ := tank.buckets["test"]

	if bucket.Rest() != 0 {

		t.Error("TokenBucket Fire() method not valid with Tank implementation of Check method")
	}

	tank.Flush()

	if bucket.Rest() != uint32(1) {

		t.Error("TokenBucket Update() method valid with Tank implementation of Flush() method")
	}
}

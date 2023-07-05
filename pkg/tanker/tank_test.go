package tanker

import (
	"testing"
)

func TestNewTank(t *testing.T) {

	tank := NewTank(NewBucketStub)

	if tank.buckets == nil {

		t.Error("fail init tanker ip_tree storage")
	}

	if tank.bucketBuilder == nil {

		t.Error("fail init tanker builder")
	}
}

func TestTank_Add(t *testing.T) {

	tank := NewTank(NewBucketStub)

	bucket := tank.Add("test")

	if _, ok := tank.buckets["test"]; !ok {

		t.Error("Tank Add() method failed")
	}

	if bucket == nil {

		t.Error("Tank Add() method return ip_tree failed")
	}

	if buf, _ := tank.buckets["test"]; bucket != buf {

		t.Error("Tank Add() method failed buckets handler not equal")
	}
}

func TestTank_Check(t *testing.T) {

	tank := NewTank(NewBucketStub)

	if !tank.Check("test") {

		t.Error("Tank Check method failed with new ip_tree")
	}

	if tank.Check("test") {

		t.Error("Tank Check method failed with existed ip_tree")
	}
}

func TestTank_Flush(t *testing.T) {

	tank := NewTank(NewBucketStub)

	_ = tank.Check("test")

	tank.Flush()

	if _, ok := tank.buckets["test"]; ok {

		t.Error("Tank Flush() method failed")
	}
}

func TestTank_Remove(t *testing.T) {

	tank := NewTank(NewBucketStub)

	_ = tank.Check("test")

	tank.Remove("test")

	if _, ok := tank.buckets["test"]; ok {

		t.Error("Tank Remove() method failed")
	}
}

type BucketStub struct {
	state bool
}

func NewBucketStub() Bucket {

	return &BucketStub{state: true}
}

func (b *BucketStub) Fire() bool {

	if b.state {
		b.state = false

		return true
	}

	return false
}

func (b *BucketStub) Update() uint32 {
	b.state = false

	return uint32(0)
}

func (b *BucketStub) Rest() uint32 {

	if b.state {
		return uint32(1)
	}

	return uint32(0)
}

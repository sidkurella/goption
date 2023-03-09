package iterator_test

import (
	"testing"

	"github.com/sidkurella/goption/iterator"
	"github.com/sidkurella/goption/option"
)

type afterNothingIterator struct {
	data []int
	i    int
}

func (a *afterNothingIterator) Next() option.Option[int] {
	if a.i >= len(a.data) {
		a.i = 0 // Loop after Nothing.
		return option.Nothing[int]()
	}
	ret := option.Some(a.data[a.i])
	a.i++
	return ret
}

func TestFuse(t *testing.T) {
	iter := &afterNothingIterator{
		data: []int{1, 2},
	}
	if iter.Next().Unwrap() != 1 {
		t.Fail()
	}
	if iter.Next().Unwrap() != 2 {
		t.Fail()
	}
	if !iter.Next().IsNothing() {
		t.Fail()
	}
	if iter.Next().Unwrap() != 1 {
		t.Fail()
	}

	fuseInner := &afterNothingIterator{
		data: []int{1, 2},
	}
	fused := iterator.Fuse[int](fuseInner)

	if fused.Next().Unwrap() != 1 {
		t.Fail()
	}
	if fused.Next().Unwrap() != 2 {
		t.Fail()
	}
	if !fused.Next().IsNothing() {
		t.Fail()
	}
	if !fused.Next().IsNothing() {
		t.Fail()
	}
	if !fused.Next().IsNothing() {
		t.Fail()
	}
}

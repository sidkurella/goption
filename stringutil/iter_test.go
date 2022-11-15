package stringutil_test

import (
	"testing"

	"github.com/sidkurella/goption/stringutil"
)

func TestIter(t *testing.T) {
	t.Run("ascii", func(t *testing.T) {
		iter := stringutil.Iter("123")
		if iter.Next().Unwrap() != '1' {
			t.Fail()
		}
		if iter.Next().Unwrap() != '2' {
			t.Fail()
		}
		if iter.Next().Unwrap() != '3' {
			t.Fail()
		}
		if !iter.Next().IsNothing() {
			t.Fail()
		}
	})
	t.Run("UTF-8", func(t *testing.T) {
		iter := stringutil.Iter("1£€3")
		if iter.Next().Unwrap() != '1' {
			t.Fail()
		}
		if iter.Next().Unwrap() != '£' {
			t.Fail()
		}
		if iter.Next().Unwrap() != '€' {
			t.Fail()
		}
		if iter.Next().Unwrap() != '3' {
			t.Fail()
		}
		if !iter.Next().IsNothing() {
			t.Fail()
		}
	})
}

func TestByteIter(t *testing.T) {
	t.Run("ascii", func(t *testing.T) {
		iter := stringutil.ByteIter("123")
		if iter.Next().Unwrap() != '1' {
			t.Fail()
		}
		if iter.Next().Unwrap() != '2' {
			t.Fail()
		}
		if iter.Next().Unwrap() != '3' {
			t.Fail()
		}
		if !iter.Next().IsNothing() {
			t.Fail()
		}
	})
	t.Run("UTF-8", func(t *testing.T) {
		iter := stringutil.ByteIter("1£€3")
		if iter.Next().Unwrap() != '1' {
			t.Fail()
		}
		if iter.Next().Unwrap() != 0xC2 {
			t.Fail()
		}
		if iter.Next().Unwrap() != 0xA3 {
			t.Fail()
		}
		if iter.Next().Unwrap() != 0xE2 {
			t.Fail()
		}
		if iter.Next().Unwrap() != 0x82 {
			t.Fail()
		}
		if iter.Next().Unwrap() != 0xAC {
			t.Fail()
		}
		if iter.Next().Unwrap() != '3' {
			t.Fail()
		}
		if !iter.Next().IsNothing() {
			t.Fail()
		}
	})
}

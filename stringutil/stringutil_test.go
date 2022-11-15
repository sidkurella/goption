package stringutil_test

import (
	"testing"

	"github.com/sidkurella/goption/stringutil"
)

func TestTruncate(t *testing.T) {
	t.Run("ascii, length less than max", func(t *testing.T) {
		expected := "123"
		actual := stringutil.Truncate("123", 100)
		if expected != actual {
			t.Fail()
		}
	})
	t.Run("ascii, length greater than max", func(t *testing.T) {
		expected := "123"
		actual := stringutil.Truncate("1234", 3)
		if expected != actual {
			t.Fail()
		}
	})
	t.Run("UTF-8", func(t *testing.T) {
		s := "1£€3"
		expected := "1£€"
		actual := stringutil.Truncate(s, 3)
		if expected != actual {
			t.Fail()
		}
	})
}

func TestTruncateBytes(t *testing.T) {
	t.Run("ascii, length less than max", func(t *testing.T) {
		expected := "123"
		actual := stringutil.TruncateBytes("123", 100)
		if expected != actual {
			t.Fail()
		}
	})
	t.Run("ascii, length greater than max", func(t *testing.T) {
		expected := "123"
		actual := stringutil.TruncateBytes("1234", 3)
		if expected != actual {
			t.Fail()
		}
	})
	t.Run("UTF-8", func(t *testing.T) {
		s := "1£€3"
		expected := "1£"
		actual := stringutil.TruncateBytes(s, 3)
		if expected != actual {
			t.Fail()
		}
	})
}

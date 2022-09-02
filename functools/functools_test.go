package functools_test

import (
	"strconv"
	"testing"

	"github.com/sidkurella/goption/functools"
)

func TestCurry(t *testing.T) {
	calls := 0
	f := func(i int, i2 int64) string {
		calls++
		return strconv.FormatInt(int64(i)*i2, 10)
	}
	curried := functools.Curry(f)
	f3 := curried(3)
	if calls != 0 { // Shouldn't be called yet.
		t.FailNow()
	}

	out := f3(8)
	expected := "24"
	if out != expected {
		t.Fail()
	}
}

func TestUncurry(t *testing.T) {
	outerCalls := 0
	innerCalls := 0
	f := func(i int) func(int64) string {
		outerCalls++
		return func(i2 int64) string {
			innerCalls++
			return strconv.FormatInt(int64(i)*i2, 10)
		}
	}
	uncurried := functools.Uncurry(f)
	out := uncurried(3, 8)
	expected := "24"
	if out != expected || outerCalls != 1 || innerCalls != 1 {
		t.Fail()
	}
}

func TestCompose(t *testing.T) {
	fCalls := 0
	gCalls := 0
	f := func(s string) float64 {
		fCalls++
		return float64(len(s) * 2)
	}
	g := func(i int) string {
		gCalls++
		return strconv.Itoa(i + 10)
	}
	composed := functools.Compose(f, g)
	out := composed(1)
	expected := float64(4)
	if out != expected || fCalls != 1 || gCalls != 1 {
		t.Fail()
	}
}

func TestMemoize(t *testing.T) {
	fCalls := 0
	f := func(i int) string {
		fCalls++
		return strconv.Itoa(i)
	}
	memoized := functools.Memoize(f)

	out := memoized(100)
	if out != "100" || fCalls != 1 {
		t.Fail()
	}
	outAgain := memoized(100)
	if out != outAgain || fCalls != 1 { // Shouldn't result in a function call.
		t.Fail()
	}
	out = memoized(200)
	if out != "200" || fCalls != 2 { // New call should be computed.
		t.Fail()
	}
	outAgain = memoized(200)
	if out != outAgain || fCalls != 2 { // Shouldn't result in a function call.
		t.Fail()
	}
}

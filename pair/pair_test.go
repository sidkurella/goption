package pair_test

import (
	"reflect"
	"testing"

	"github.com/sidkurella/goption/pair"
)

func TestPair_From(t *testing.T) {
	expected := pair.Pair[int, string]{
		First:  10,
		Second: "300",
	}
	actual := pair.From(10, "300")
	if !reflect.DeepEqual(actual, expected) {
		t.Fail()
	}
}

func TestPair_Into(t *testing.T) {
	p := pair.Pair[int, string]{
		First:  10,
		Second: "300",
	}
	fst, snd := p.Into()
	if !reflect.DeepEqual(fst, 10) {
		t.Fail()
	}
	if !reflect.DeepEqual(snd, "300") {
		t.Fail()
	}
}

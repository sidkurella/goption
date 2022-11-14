package iterator

import (
	"github.com/sidkurella/goption/option"
	"golang.org/x/exp/constraints"
)

type numeric interface {
	constraints.Integer | constraints.Float
}

type rangeIterator[T numeric] struct {
	start      T
	end        T
	current    T
	step       T
	includeEnd bool
	backwards  bool
}

func (r *rangeIterator[T]) Next() option.Option[T] {
	var ret option.Option[T] = option.Nothing[T]{}
	if (!r.backwards && r.current < r.end) ||
		(r.backwards && r.current > r.end) ||
		(r.current == r.end && r.includeEnd) {
		ret = option.Some[T]{Value: r.current}
		r.current += r.step
	}

	return ret
}

// Returns an iterator ranging from start (inclusive) to end (exclusive), stepping by 1.
// If end is less than start, the iterator will be empty.
func Range[T numeric](start T, end T) *rangeIterator[T] {
	return &rangeIterator[T]{
		start:      start,
		end:        end,
		current:    start,
		step:       T(1),
		includeEnd: false,
		backwards:  false,
	}
}

// Returns an iterator ranging from start (inclusive) to end (exclusive), stepping by step.
// If end is less than start, the iterator will be empty.
// This can be used with a negative step. If so, if end is greater than start, the iterator will be empty.
// NOTE: A zero step will return start ad infinitum.
func RangeBy[T numeric](start T, end T, step T) *rangeIterator[T] {
	return &rangeIterator[T]{
		start:      start,
		end:        end,
		current:    start,
		step:       step,
		includeEnd: false,
		backwards:  step < 0,
	}
}

// Returns an iterator ranging from start (inclusive) to end (inclusive), stepping by 1.
// If end is less than start, the iterator will be empty.
func RangeInclusive[T numeric](start T, end T) *rangeIterator[T] {
	return &rangeIterator[T]{
		start:      start,
		end:        end,
		current:    start,
		step:       T(1),
		includeEnd: true,
		backwards:  false,
	}
}

// Returns an iterator ranging from start (inclusive) to end (exclusive), stepping by step.
// If end is less than start, the iterator will be empty.
// This can be used with a negative step. If so, if end is greater than start, the iterator will be empty.
// NOTE: A zero step will return start ad infinitum.
func RangeInclusiveBy[T numeric](start T, end T, step T) *rangeIterator[T] {
	return &rangeIterator[T]{
		start:      start,
		end:        end,
		current:    start,
		step:       step,
		includeEnd: true,
		backwards:  step < 0,
	}
}

# goption

`goption` is a small utility library that adds Rust/FP-inspired types and helpers to Go.

It includes:
- `Option`, `Result`, and `Either`
- iterator adapters and collectors
- `Set` and map helpers
- slice and function utility helpers

## Installation

```bash
go get github.com/sidkurella/goption@latest
```

## Packages

- `option`: optional values (`Some` / `Nothing`), plus JSON and SQL compatibility helpers
- `result`: success/error values (`Ok` / `Err`)
- `either`: two-branch values (`First` / `Second`)
- `iterator`: pull-based iterator adapters and collectors
- `set`: hash set utilities and set algebra
- `maputil`: map wrappers and transforms
- `sliceutil`: common slice transforms and prefix/suffix helpers
- `functools`: curry/uncurry/compose/memoize helpers

## Option Basics

```go
package main

import (
	"fmt"

	"github.com/sidkurella/goption/option"
)

func main() {
	name := option.Some("sid")
	none := option.Nothing[string]()

	upper := option.Map(name, func(s string) string { return s + "!" })
	fallback := none.UnwrapOr("guest")

	fmt.Println(upper.Unwrap()) // sid!
	fmt.Println(fallback)       // guest
}
```

## Result and Either Basics

```go
package main

import (
	"fmt"
	"strconv"

	"github.com/sidkurella/goption/either"
	"github.com/sidkurella/goption/result"
)

func parseInt(s string) result.Result[int, error] {
	v, err := strconv.Atoi(s)
	return result.From(v, err)
}

func main() {
	r := result.AndThen(parseInt("42"), func(v int) result.Result[int, error] {
		return result.Ok[int, error](v * 2)
	})
	fmt.Println(r) // Ok(84)

	e := either.First[int, string](10)
	fmt.Println(e.IsFirst()) // true
}
```

## Iterators

The iterator package supports adapter chains and collectors similar to Rust-style iteration.

```go
package main

import (
	"fmt"

	"github.com/sidkurella/goption/iterator"
	"github.com/sidkurella/goption/sliceutil"
)

func main() {
	it := sliceutil.Iter([]int{1, 2, 3, 4, 5, 6})
	mapped := iterator.Map(it, func(v int) int { return v * 10 })
	filtered := iterator.Filter(mapped, func(v int) bool { return v >= 30 })
	firstTwo := iterator.Take(filtered, 2)

	out := iterator.Collect(firstTwo)
	fmt.Println(out) // [30 40]
}
```

### `iter.Seq` / `iter.Seq2` interop

`FromSeq` and `FromSeq2` use `iter.Pull`/`iter.Pull2` internally and hold cleanup resources. If you may not fully exhaust the iterator, defer `Close`.

```go
package main

import (
	"fmt"
	"slices"

	"github.com/sidkurella/goption/iterator"
)

func main() {
	it := iterator.FromSeq(slices.Values([]int{1, 2, 3, 4, 5}))
	defer it.Close()

	first := iterator.Take(it, 2)
	fmt.Println(iterator.Collect(first)) // [1 2]
}
```

## Set Utilities

```go
package main

import (
	"fmt"

	"github.com/sidkurella/goption/set"
)

func main() {
	a := set.FromSlice([]int{1, 2, 3})
	b := set.FromSlice([]int{3, 4, 5})

	intersection := a.Intersection(b)
	pairDiff := a.PairedDifference(b)

	fmt.Println(intersection.Contains(3))    // true
	fmt.Println(pairDiff.First.Contains(1))  // true (only in a)
	fmt.Println(pairDiff.Second.Contains(5)) // true (only in b)
}
```

## Map and Slice Helpers

```go
package main

import (
	"fmt"

	"github.com/sidkurella/goption/maputil"
	"github.com/sidkurella/goption/sliceutil"
)

func main() {
	m := maputil.New[string, int]()
	m.Insert("a", 1)
	m.Insert("b", 2)

	doubled := maputil.Apply(m, func(k string, v int) (string, int) {
		return k, v * 2
	})
	fmt.Println(doubled.Get("b").Unwrap()) // 4

	out := sliceutil.Map([]int{1, 2, 3}, func(v int) int { return v * 3 })
	fmt.Println(out) // [3 6 9]
}
```

## Best Practices

- Use `Option` when absence is expected and not an error condition.
- Use `Result` when you need to preserve an error payload.
- Prefer iterator adapters for composable pipelines; call `Collect` at the boundary.
- When using `FromSeq`/`FromSeq2`, `defer it.Close()` right after creation if full exhaustion is not guaranteed.
- Set and map iteration order is not stable; sort collected outputs in tests if order matters.

## Common Pitfalls

- Missing JSON field vs `null` for `Option`: `null` becomes `Nothing`; a missing field keeps the existing field value unless decoding into a fresh zero-value struct.
- Calling `Unwrap` on `Nothing`/`Err`/`Second` panics; use `UnwrapOr`, `Match`, or conversion helpers when uncertain.
- Assuming iterator adapters auto-close upstream pull iterators; they do not.

## Development

Run all tests:

```bash
go test ./...
```

# go-toolkit

Generic helpers for slices, maps, optionals, pointers, and concurrency primitives in Go.

Small, composable building blocks for common collection workflows. Focused on
explicit semantics and predictable behavior over clever abstractions.

## Install

```bash
go get github.com/zaurkadiev/go-toolkit
```

Requires Go 1.25+.

## Packages

| Package | Purpose |
|---------|---------|
| [`tslices`](./tslices) | Map, filter, group, dedup, aggregate, search over slices |
| [`tmaps`](./tmaps) | Map ops, filter, merge, convert to slices |
| [`optional`](./optional) | Typed optional for absent/zero values |
| [`pointers`](./pointers) | Pointer construction and nil-safe dereference |
| [`future`](./future) | Run a function in the background and await its result |
| [`tsync`](./tsync) | Typed `Map`, `RWMap` and a context-aware `Semaphore` |
| [`worker_pool`](./worker_pool) | Run jobs on a fixed number of workers, stop on first error |

## Examples

### tslices

```go
import "github.com/zaurkadiev/go-toolkit/tslices"

type User struct {
    ID   int
    Name string
}

users := []User{{1, "a"}, {2, "b"}, {3, "a"}}

names := tslices.Map(users, func(u User) (string, bool) {
    return u.Name, true
})
// ["a", "b", "a"]

byName := tslices.GroupBy(users, func(u User) string { return u.Name })
// map[a:[{1 a} {3 a}] b:[{2 b}]]
```

### tmaps

```go
import "github.com/zaurkadiev/go-toolkit/tmaps"

counts := map[string]int{"a": 1, "b": 2, "c": 3}

keys := tmaps.GetKeysSorted(counts)
// ["a", "b", "c"]

filtered := tmaps.Filter(counts, func(k string, v int) bool {
    return v > 1
})
// map[b:2 c:3]
```

### optional

```go
import "github.com/zaurkadiev/go-toolkit/optional"

opt := optional.Of(42)
if v, ok := opt.Get(); ok {
    fmt.Println(v) // 42
}

empty := optional.Empty[int]()
fmt.Println(empty.OrElse(0)) // 0
```

### pointers

```go
import "github.com/zaurkadiev/go-toolkit/pointers"

p := pointers.To("hello")    // *string
v := pointers.Deref(p, "")   // "hello"
v2 := pointers.Deref[string](nil, "fallback") // "fallback"
```

## Design

- Mapper callbacks return `(value, ok)` so filtering is built in.
- Error-aware variants (`*Err`) short-circuit on first error.
- Map iteration order is non-deterministic — helpers that return "first" do not
  guarantee ordering.
- Narrow type constraints: `comparable` for keys, `cmp.Ordered` for sort paths,
  `any` only when no constraint applies.

## License

MIT — see [LICENSE](./LICENSE).

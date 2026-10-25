# go-toolkit

Generic helpers for slices, maps, optionals, pointers, and concurrency primitives in Go.

Small, composable building blocks for common collection workflows. Focused on
explicit semantics and predictable behavior over clever abstractions.

## Install

```bash
go get github.com/zaurkadiev/go-toolkit
```

Requires Go 1.27+.

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

byName := tslices.Group(users, func(u User) (string, User, bool) { return u.Name, u, true })
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

n := 42
opt := optional.New(&n)
fmt.Println(opt.Get()) // 42

empty := optional.Empty[int]()
fmt.Println(empty.IsEmpty())       // true
fmt.Println(empty.GetOrDefault(7)) // 7
```

### pointers

```go
import "github.com/zaurkadiev/go-toolkit/pointers"

p := pointers.Of("hello")           // *string
v := pointers.ValueOf(p)            // "hello"
v2 := pointers.ValueOf[string](nil) // "" (zero value)
```

### future

```go
import "github.com/zaurkadiev/go-toolkit/future"

f := future.New(func() (int, error) { return 42, nil })
v, err := f.Await() // 42, nil
```

### tsync

```go
import "github.com/zaurkadiev/go-toolkit/tsync"

var m tsync.Map[string, int]
m.Store("a", 1)
v, ok := m.Load("a") // 1, true

sem := tsync.NewSemaphore(2)
if err := sem.Acquire(ctx); err == nil {
    defer sem.Release()
}
```

### worker_pool

```go
import "github.com/zaurkadiev/go-toolkit/worker_pool"

wp := worker_pool.New[int](4)
wp.Add(func() (int, error) { return 1, nil })
wp.Add(func() (int, error) { return 2, nil })
res, err := wp.Run(ctx) // [1 2], nil (results in job order)
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

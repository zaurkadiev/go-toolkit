package tsync

import "sync"

// Map is a typed wrapper over sync.Map. Its methods mirror sync.Map with typed
// keys and values.
//
// Map is optimized for concurrent access with low lock contention:
//   - write-once, read-many entries, such as caches that only grow;
//   - goroutines that read, write and overwrite disjoint sets of keys;
//   - highly contended maps where a single mutex becomes a bottleneck.
//
// Prefer RWMap, a map guarded by sync.RWMutex, when:
//   - several keys must change together or invariants span multiple keys,
//     since Map makes each operation atomic only on its own;
//   - the map length is needed, since Map has no Len;
//   - a consistent snapshot is needed, since Range does not provide one;
//   - contention is low or the map is used from a single goroutine;
//   - allocations matter, since values are stored as any and non-pointer
//     values usually allocate.
//
// The zero Map is empty and ready for use. A Map must not be copied after first use.
//
// Example:
//
//	var m Map[string, int]
//	m.Store("a", 1)
//	v, ok := m.Load("a")
//	// v == 1, ok == true
type Map[K comparable, V any] struct {
	m sync.Map
}

// Load returns the value stored in the map for a key, or the zero value of V if no
// value is present. The ok result indicates whether value was found in the map.
//
// Example:
//
//	var m Map[string, int]
//	m.Store("a", 1)
//	v, ok := m.Load("a")
//	// v == 1, ok == true
//	v, ok = m.Load("b")
//	// v == 0, ok == false
func (m *Map[K, V]) Load(key K) (value V, ok bool) {
	v, ok := m.m.Load(key)

	return cast[V](v), ok
}

// Store sets the value for a key.
//
// Example:
//
//	var m Map[string, int]
//	m.Store("a", 1)
//	m.Store("a", 2)
//	v, _ := m.Load("a")
//	// v == 2
func (m *Map[K, V]) Store(key K, value V) {
	m.m.Store(key, value)
}

// LoadOrStore returns the existing value for the key if present.
// Otherwise, it stores and returns the given value.
// The loaded result is true if the value was loaded, false if stored.
//
// Example:
//
//	var m Map[string, int]
//	v, loaded := m.LoadOrStore("a", 1)
//	// v == 1, loaded == false
//	v, loaded = m.LoadOrStore("a", 2)
//	// v == 1, loaded == true
func (m *Map[K, V]) LoadOrStore(key K, value V) (actual V, loaded bool) {
	v, loaded := m.m.LoadOrStore(key, value)

	return cast[V](v), loaded
}

// LoadAndDelete deletes the value for a key, returning the previous value if any.
// The loaded result reports whether the key was present.
//
// Example:
//
//	var m Map[string, int]
//	m.Store("a", 1)
//	v, loaded := m.LoadAndDelete("a")
//	// v == 1, loaded == true
//	v, loaded = m.LoadAndDelete("a")
//	// v == 0, loaded == false
func (m *Map[K, V]) LoadAndDelete(key K) (value V, loaded bool) {
	v, loaded := m.m.LoadAndDelete(key)

	return cast[V](v), loaded
}

// Delete deletes the value for a key.
//
// Example:
//
//	var m Map[string, int]
//	m.Store("a", 1)
//	m.Delete("a")
//	_, ok := m.Load("a")
//	// ok == false
func (m *Map[K, V]) Delete(key K) {
	m.m.Delete(key)
}

// Swap swaps the value for a key and returns the previous value if any.
// The loaded result reports whether the key was present.
//
// Example:
//
//	var m Map[string, int]
//	prev, loaded := m.Swap("a", 1)
//	// prev == 0, loaded == false
//	prev, loaded = m.Swap("a", 2)
//	// prev == 1, loaded == true
func (m *Map[K, V]) Swap(key K, value V) (previous V, loaded bool) {
	v, loaded := m.m.Swap(key, value)

	return cast[V](v), loaded
}

// CompareAndSwap swaps the old and new values for key
// if the value stored in the map is equal to old.
//
// As with sync.Map, it panics if V is not a comparable type.
//
// Example:
//
//	var m Map[string, int]
//	m.Store("a", 1)
//	swapped := m.CompareAndSwap("a", 1, 2)
//	// swapped == true, m.Load("a") == 2
//	swapped = m.CompareAndSwap("a", 1, 3)
//	// swapped == false, value stays 2
func (m *Map[K, V]) CompareAndSwap(key K, old, new V) (swapped bool) {
	return m.m.CompareAndSwap(key, old, new)
}

// CompareAndDelete deletes the entry for key if its value is equal to old.
// If there is no current value for key in the map, CompareAndDelete returns false.
//
// As with sync.Map, it panics if V is not a comparable type.
//
// Example:
//
//	var m Map[string, int]
//	m.Store("a", 1)
//	deleted := m.CompareAndDelete("a", 2)
//	// deleted == false
//	deleted = m.CompareAndDelete("a", 1)
//	// deleted == true
func (m *Map[K, V]) CompareAndDelete(key K, old V) (deleted bool) {
	return m.m.CompareAndDelete(key, old)
}

// Range calls f sequentially for each key and value present in the map.
// If f returns false, range stops the iteration. Range has the same consistency
// guarantees as sync.Map.Range.
//
// Example:
//
//	var m Map[string, int]
//	m.Store("a", 1)
//	m.Store("b", 2)
//	sum := 0
//	m.Range(func(key string, value int) bool {
//		sum += value
//		return true // return false to stop
//	})
//	// sum == 3
func (m *Map[K, V]) Range(f func(key K, value V) bool) {
	m.m.Range(func(k, v any) bool {
		return f(cast[K](k), cast[V](v))
	})
}

// Clear deletes all the entries, resulting in an empty Map.
//
// Example:
//
//	var m Map[string, int]
//	m.Store("a", 1)
//	m.Clear()
//	_, ok := m.Load("a")
//	// ok == false
func (m *Map[K, V]) Clear() {
	m.m.Clear()
}

// cast converts a stored value back to T. A nil value (absent key, or a nil
// interface stored for an interface type T) yields the zero value of T instead
// of panicking.
func cast[T any](v any) T {
	t, _ := v.(T)

	return t
}

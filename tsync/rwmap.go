package tsync

import (
	"maps"
	"sync"
)

// RWMap is a typed map guarded by sync.RWMutex. Its methods mirror Map, plus Len.
//
// Prefer RWMap over Map when:
//   - the map length is needed;
//   - a consistent snapshot is needed during iteration;
//   - contention is low or writes are rare;
//   - allocations matter, since values are stored without boxing to any.
//
// Each method is atomic on its own; a sequence of calls is not.
//
// The zero RWMap is empty and ready for use. An RWMap must not be copied after
// first use.
//
// Example:
//
//	var m RWMap[string, int]
//	m.Store("a", 1)
//	v, ok := m.Load("a")
//	// v == 1, ok == true
type RWMap[K comparable, V any] struct {
	mu sync.RWMutex
	m  map[K]V
}

// Load returns the value stored in the map for a key, or the zero value of V if no
// value is present. The ok result indicates whether value was found in the map.
//
// Example:
//
//	var m RWMap[string, int]
//	m.Store("a", 1)
//	v, ok := m.Load("a")
//	// v == 1, ok == true
//	v, ok = m.Load("b")
//	// v == 0, ok == false
func (m *RWMap[K, V]) Load(key K) (value V, ok bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	value, ok = m.m[key]

	return value, ok
}

// Store sets the value for a key.
//
// Example:
//
//	var m RWMap[string, int]
//	m.Store("a", 1)
//	m.Store("a", 2)
//	v, _ := m.Load("a")
//	// v == 2
func (m *RWMap[K, V]) Store(key K, value V) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.init()
	m.m[key] = value
}

// LoadOrStore returns the existing value for the key if present.
// Otherwise, it stores and returns the given value.
// The loaded result is true if the value was loaded, false if stored.
//
// Example:
//
//	var m RWMap[string, int]
//	v, loaded := m.LoadOrStore("a", 1)
//	// v == 1, loaded == false
//	v, loaded = m.LoadOrStore("a", 2)
//	// v == 1, loaded == true
func (m *RWMap[K, V]) LoadOrStore(key K, value V) (actual V, loaded bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if actual, loaded = m.m[key]; loaded {
		return actual, true
	}
	m.init()
	m.m[key] = value

	return value, false
}

// LoadAndDelete deletes the value for a key, returning the previous value if any.
// The loaded result reports whether the key was present.
//
// Example:
//
//	var m RWMap[string, int]
//	m.Store("a", 1)
//	v, loaded := m.LoadAndDelete("a")
//	// v == 1, loaded == true
//	v, loaded = m.LoadAndDelete("a")
//	// v == 0, loaded == false
func (m *RWMap[K, V]) LoadAndDelete(key K) (value V, loaded bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	value, loaded = m.m[key]
	delete(m.m, key)

	return value, loaded
}

// Delete deletes the value for a key.
//
// Example:
//
//	var m RWMap[string, int]
//	m.Store("a", 1)
//	m.Delete("a")
//	_, ok := m.Load("a")
//	// ok == false
func (m *RWMap[K, V]) Delete(key K) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.m, key)
}

// Swap swaps the value for a key and returns the previous value if any.
// The loaded result reports whether the key was present.
//
// Example:
//
//	var m RWMap[string, int]
//	prev, loaded := m.Swap("a", 1)
//	// prev == 0, loaded == false
//	prev, loaded = m.Swap("a", 2)
//	// prev == 1, loaded == true
func (m *RWMap[K, V]) Swap(key K, value V) (previous V, loaded bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	previous, loaded = m.m[key]
	m.init()
	m.m[key] = value

	return previous, loaded
}

// CompareAndSwap swaps the old and new values for key
// if the value stored in the map is equal to old.
// If there is no current value for key in the map, CompareAndSwap returns false.
//
// As with Map, it panics if V is not a comparable type.
//
// Example:
//
//	var m RWMap[string, int]
//	m.Store("a", 1)
//	swapped := m.CompareAndSwap("a", 1, 2)
//	// swapped == true, m.Load("a") == 2
//	swapped = m.CompareAndSwap("a", 1, 3)
//	// swapped == false, value stays 2
func (m *RWMap[K, V]) CompareAndSwap(key K, old, new V) (swapped bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	cur, ok := m.m[key]
	if !ok || any(cur) != any(old) {
		return false
	}
	m.m[key] = new

	return true
}

// CompareAndDelete deletes the entry for key if its value is equal to old.
// If there is no current value for key in the map, CompareAndDelete returns false.
//
// As with Map, it panics if V is not a comparable type.
//
// Example:
//
//	var m RWMap[string, int]
//	m.Store("a", 1)
//	deleted := m.CompareAndDelete("a", 2)
//	// deleted == false
//	deleted = m.CompareAndDelete("a", 1)
//	// deleted == true
func (m *RWMap[K, V]) CompareAndDelete(key K, old V) (deleted bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	cur, ok := m.m[key]
	if !ok || any(cur) != any(old) {
		return false
	}
	delete(m.m, key)

	return true
}

// Range calls f sequentially for each key and value present in the map.
// If f returns false, range stops the iteration.
//
// Range iterates over a snapshot taken under the read lock, so f sees a consistent
// view of the map and may call other RWMap methods without deadlocking. Changes
// made during iteration are not visible to it.
//
// Example:
//
//	var m RWMap[string, int]
//	m.Store("a", 1)
//	m.Store("b", 2)
//	sum := 0
//	m.Range(func(key string, value int) bool {
//		sum += value
//		return true // return false to stop
//	})
//	// sum == 3
func (m *RWMap[K, V]) Range(f func(key K, value V) bool) {
	m.mu.RLock()
	snapshot := maps.Clone(m.m)
	m.mu.RUnlock()

	for k, v := range snapshot {
		if !f(k, v) {
			return
		}
	}
}

// Clear deletes all the entries, resulting in an empty RWMap.
//
// Example:
//
//	var m RWMap[string, int]
//	m.Store("a", 1)
//	m.Clear()
//	_, ok := m.Load("a")
//	// ok == false
func (m *RWMap[K, V]) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()
	clear(m.m)
}

// Len returns the number of entries in the map.
//
// Example:
//
//	var m RWMap[string, int]
//	m.Store("a", 1)
//	m.Store("b", 2)
//	n := m.Len()
//	// n == 2
func (m *RWMap[K, V]) Len() int {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return len(m.m)
}

// init allocates the underlying map on first write. Callers must hold the write lock.
func (m *RWMap[K, V]) init() {
	if m.m == nil {
		m.m = make(map[K]V)
	}
}

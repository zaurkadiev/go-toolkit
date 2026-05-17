package tmaps

// Merge combines multiple maps into a single map.
// When keys overlap, values from later maps overwrite earlier ones.
//
// Behavior details:
// - With a nil variadic input, it returns nil.
// - With an empty but non-nil input slice, it returns an empty map.
// - If all input maps are empty, it returns an empty map.
// - The output map has no guaranteed ordering.
//
// Example:
//
//	out := Merge(
//		map[int]string{1: "a", 2: "b"},
//		map[int]string{2: "bb", 3: "c"},
//	)
//	// out == map[int]string{1: "a", 2: "bb", 3: "c"}
func Merge[K comparable, V any](maps ...map[K]V) map[K]V {
	if maps == nil {
		return nil
	}
	if len(maps) == 0 {
		return map[K]V{}
	}

	c := 0
	for _, m := range maps {
		c += len(m)
	}
	if c == 0 {
		return map[K]V{}
	}

	r := make(map[K]V, c)
	for _, m := range maps {
		for k, v := range m {
			r[k] = v
		}
	}

	return r
}

package tslices

type integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

type float interface {
	~float32 | ~float64
}

type ordered interface {
	float | integer | ~string
}

type Pair[L, R any] struct {
	Left  L
	Right R
}

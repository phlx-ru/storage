package slices

type Scalar interface {
	~bool |
		~string |
		~int8 | ~int16 | ~int32 | ~int64 | ~int |
		~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uint |
		~float32 | ~float64 |
		~complex64 | ~complex128
}

func Includes[T Scalar](needle T, haystack []T) bool {
	for _, current := range haystack {
		if current == needle {
			return true
		}
	}
	return false
}

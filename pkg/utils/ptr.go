package utils

//go:fix inline
func PtrOf[T any](v T) *T {
	return new(v)
}

func DerefOrZero[T any](ptr *T) T {
	if ptr == nil {
		var zero T
		return zero
	}
	return *ptr
}

package utils

func PtrOf[T any](v T) *T {
	return &v
}

func DerefOrZero[T any](ptr *T) T {
	if ptr == nil {
		var zero T
		return zero
	}
	return *ptr
}

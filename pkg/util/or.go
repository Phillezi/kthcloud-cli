package util

// Returns the first non default value, fallback is default value
func Or[T comparable](vals ...T) T {
	var zero T
	for _, v := range vals {
		if v != zero {
			return v
		}
	}
	return zero
}

// Like Or but the first argument is a pointer, that if it is not nil is used
func PtrOr[T comparable](ptr *T, vals ...T) T {
	if ptr != nil {
		return *ptr
	}
	return Or(vals...)
}

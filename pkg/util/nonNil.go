package util

func IfNotNilDo[T any](e *T, action func()) {
	if e != nil {
		action()
	}
}

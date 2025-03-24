package util

func ToPointer[T any](value T) *T {
	return &value
}

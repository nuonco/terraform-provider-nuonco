package provider

// toPtr generically returns a reference to the value v of type T
func toPtr[T any](v T) *T {
	return &v
}

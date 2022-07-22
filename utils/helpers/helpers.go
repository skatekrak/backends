package helpers

func SetIfNotNil[T any](optionalValue *T, value T) T {
	if optionalValue != nil {
		return *optionalValue
	}
	return value
}

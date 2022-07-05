package helpers

func Has[T comparable](slice []T, element T) bool {
	for _, el := range slice {
		if element == el {
			return true
		}
	}
	return false
}

func Find[T any](m []*T, iterate func(*T) bool) (*T, bool) {
	for _, item := range m {
		result := iterate(item)
		if result {
			return item, true
		}
	}

	var empty T
	return &empty, false
}

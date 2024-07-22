package utils

func IsExist[T any](slice []T, item T, eq func(T, T) bool) bool {
	for _, v := range slice {
		if eq(v, item) {
			return true
		}
	}
	return false
}

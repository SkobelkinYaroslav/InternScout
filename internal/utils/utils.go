package utils

func IsExist[T any](slice []T, item T, eq func(T, T) bool) bool {
	for _, v := range slice {
		if eq(v, item) {
			return true
		}
	}
	return false
}

func ArrayToMapStruct[T comparable](arr []T) map[T]struct{} {
	m := make(map[T]struct{}, len(arr))
	for _, v := range arr {
		m[v] = struct{}{}
	}
	return m
}

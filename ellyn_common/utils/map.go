package utils

func GetMapKeys[K comparable, V any](m map[K]V) (res []K) {
	for k := range m {
		res = append(res, k)
	}
	return
}

func GetMapValues[K comparable, V any](m map[K]V) (res []V) {
	for _, v := range m {
		res = append(res, v)
	}
	return
}

func CopyMap[K comparable, V any](m map[K]V) (res map[K]V) {
	res = make(map[K]V)
	for k, v := range m {
		res[k] = v
	}
	return res
}

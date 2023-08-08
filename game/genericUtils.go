package game

func MapToAnyMap[K comparable, V any](m map[K]V) map[K]any {
	newMap := make(map[K]any, len(m))
	for key := range m {
		newMap[key] = true
	}
	return newMap
}

func ListToMap[I comparable](lst []I) map[I]any {
	newMap := make(map[I]any, len(lst))
	for _, val := range lst {
		newMap[val] = true
	}
	return newMap
}

func SingleToMap[I comparable](item I) map[I]any {
	newMap := make(map[I]any, 1)
	newMap[item] = true
	return newMap
}

func MapKeysToList[K comparable, V any](m map[K]V) []K {
	newList := make([]K, len(m))
	i := 0
	for key := range m {
		newList[i] = key
		i += 1
	}
	return newList
}

func MapValsToList[K comparable, V any](m map[K]V) []V {
	newList := make([]V, len(m))
	i := 0
	for _, val := range m {
		newList[i] = val
		i += 1
	}
	return newList
}

func MapValsToPointerList[K comparable, V any](m map[K]V) []*V {
	newList := make([]*V, len(m))
	i := 0
	for _, val := range m {
		newList[i] = &val
		i += 1
	}
	return newList
}

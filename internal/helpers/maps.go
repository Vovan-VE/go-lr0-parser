package helpers

import (
	"sort"
)

func MapSortedInt[K ~int, V any](m map[K]V) []MapItem[K, V] {
	return MapSorted(m, lessInt[K])
}

func MapSorted[K comparable, V any](m map[K]V, less func(a, b K) bool) []MapItem[K, V] {
	if len(m) == 0 {
		return nil
	}
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool { return less(keys[i], keys[j]) })

	res := make([]MapItem[K, V], 0, len(m))
	for _, k := range keys {
		res = append(res, MapItem[K, V]{K: k, V: m[k]})
	}
	return res
}

type MapItem[K, V any] struct {
	K K
	V V
}

func lessInt[K ~int](a, b K) bool { return a < b }

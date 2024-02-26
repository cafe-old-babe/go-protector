package utils

import (
	"go-protector/server/core/consts"
)

func SliceSub[T uint64 | int64](v1, v2 []T) (sub []T) {
	if len(v1) <= 0 {
		return
	}
	if len(v2) <= 0 {
		return v1
	}
	for i := range v1 {
		var exists bool
		for j := range v2 {
			if v1[i] == v2[j] {
				exists = true
				break
			}
		}
		if !exists {
			sub = append(sub, v1[i])
		}
	}

	return
}

func SliceSubN[T uint64 | int64](v1, v2 []T) (sub []T) {

	if len(v1) <= 0 {
		return
	}
	if len(v2) <= 0 {
		return v1
	}
	// v1-v2
	tempSet := map[T]any{}
	for _, elem := range v1 {
		tempSet[elem] = consts.EmptyVal
	}
	for _, elem := range v2 {
		delete(tempSet, elem)
	}

	for k := range tempSet {
		sub = append(sub, k)
	}
	return

}

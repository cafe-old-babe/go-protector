package utils

import (
	"go-protector/server/internal/consts"
	"reflect"
)

type name struct {
}

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

func SliceToFieldSlice[T interface{ uint64 | string }](field string, slice interface{}) (fieldSlice []T) {
	fieldSlice = make([]T, 0)
	//循环slice
	sliceValue := reflect.ValueOf(slice)
	if sliceValue.Kind() != reflect.Slice {
		return
	}

	for i := 0; i < sliceValue.Len(); i++ {

		elem := sliceValue.Index(i).Interface()

		if reflect.ValueOf(elem).Kind() != reflect.Struct {
			continue
		}
		if reflect.ValueOf(elem).FieldByName(field).IsValid() {
			continue
		}

		if val, ok := reflect.Indirect(reflect.ValueOf(elem)).FieldByName(field).Interface().(T); ok {
			fieldSlice = append(fieldSlice, val)
		}
	}

	return
}

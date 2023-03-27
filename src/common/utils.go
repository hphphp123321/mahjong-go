package common

import (
	"errors"
	"reflect"
)

func Max[T Comparable](o []T) T {
	if len(o) == 0 {
		var zero T
		return zero
	}
	maxIndex := 0
	for i := len(o) - 1; i > 0; i-- {
		if o[i].CompareTo(o[maxIndex]) > 0 {
			maxIndex = i
		}
	}
	return o[maxIndex]
}

func Min[T Comparable](o []T) T {
	if len(o) == 0 {
		var zero T
		return zero
	}
	maxIndex := 0
	for i := len(o) - 1; i > 0; i-- {
		if o[i].CompareTo(o[maxIndex]) < 0 {
			maxIndex = i
		}
	}
	return o[maxIndex]
}

func IndexOf[T Comparable](o T, t []T) int {
	if len(t) == 0 {
		return -1
	}
	for i := len(t) - 1; i >= 0; i-- {
		if t[i].CompareTo(o) == 0 {
			return i
		}
	}
	return -1
}

func MaxNum[T Num](ns []T) T {
	if len(ns) == 0 {
		var zero T
		return zero
	}
	maxIndex := 0
	for i := len(ns) - 1; i > 0; i-- {
		if ns[i]-ns[maxIndex] > 0 {
			maxIndex = i
		}
	}
	return ns[maxIndex]
}

func MinNum[T Num](ns []T) T {
	if len(ns) == 0 {
		var zero T
		return zero
	}
	minIndex := 0
	for i := len(ns) - 1; i > 0; i-- {
		if ns[i]-ns[minIndex] < 0 {
			minIndex = i
		}
	}
	return ns[minIndex]
}

func Contain(obj interface{}, target interface{}) bool {
	targetValue := reflect.ValueOf(target)
	switch reflect.TypeOf(target).Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < targetValue.Len(); i++ {
			if targetValue.Index(i).Interface() == obj {
				return true
			}
		}
	case reflect.Map:
		if targetValue.MapIndex(reflect.ValueOf(obj)).IsValid() {
			return true
		}
	}

	return false
}

func Remove(obj interface{}, target interface{}) (interface{}, error) {
	targetValue := reflect.ValueOf(target)
	switch reflect.TypeOf(target).Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < targetValue.Len(); i++ {
			if targetValue.Index(i).Interface() == obj {
				return reflect.AppendSlice(targetValue.Slice(0, i), targetValue.Slice(i+1, targetValue.Len())).Interface(), nil
			}
		}
	case reflect.Map:
		if targetValue.MapIndex(reflect.ValueOf(obj)).IsValid() {
			delete(targetValue.Interface().(map[interface{}]interface{}), obj)
			return targetValue.Interface(), nil
		}
	}
	return nil, errors.New("cannot delete: value not found")
}

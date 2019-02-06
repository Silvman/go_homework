package main

import (
	"sort"
	"strconv"
)

// сюда вам надо писать функции, которых не хватает, чтобы проходили тесты в gotchas_test.go

func ReturnInt() int {
	return 1
}

func ReturnFloat() float32 {
	return 1.1
}

func ReturnIntArray() [3]int {
	return [3]int{1, 3, 4}
}

func ReturnIntSlice() []int {
	return []int{1, 2, 3}
}

func IntSliceToString(s []int) (res string) {
	for _, value := range s {
		res += strconv.Itoa(value)
	}

	return
}

func MergeSlices(sFloat []float32, sInt []int32) []int {
	res := make([]int, len(sFloat)+len(sInt))

	for key, value := range sFloat {
		res[key] = int(value)
	}

	for key, value := range sInt {
		res[key+len(sFloat)] = int(value)
	}

	return res
}

func GetMapValuesSortedByKey(m map[int]string) []string {
	keys := make([]int, 0, len(m))
	result := make([]string, 0, len(m))

	for key := range m {
		keys = append(keys, key)
	}
	sort.Ints(keys)

	for _, value := range keys {
		result = append(result, m[value])
	}

	return result
}

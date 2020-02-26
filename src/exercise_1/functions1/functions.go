package functions1

import "fmt"
import "sort"

const startingLevel = 1

func ReturnInt() int {
	return startingLevel
}

func ReturnFloat() float32 {
	var floatVariable float32 = 1.1
    return floatVariable
}

func ReturnIntArray() [3]int {
    array := [3]int{1, 3, 4}
	return array
}

func ReturnIntSlice() []int {
    array :=  []int{1, 2, 3}
	return array
}

func IntSliceToString(array []int) string {
    a := ""
    for index := range array {
        a += fmt.Sprintf("%v",array[index])
    }

	return a
}

func MergeSlices(a []float32, b []int32) []int {
    m := make([]int, 0, len(a)+len(b))

    for _, item := range a {
        m = append(m, int(item))
    }

    for _, item := range b {
        m = append(m, int(item))
    }

    return m
}

func GetMapValuesSortedByKey(input map[int]string) []string {
	array := make([]string, len(input))
	mapKeys := getKeysFromMap(input)
	sort.Ints(mapKeys)
	for index := range mapKeys {
       array[index] = input[mapKeys[index]]
    }

    return array
}

func getKeysFromMap(input map[int]string) []int {
    keys := make([]int, len(input))
	i := 0
    for k := range input {
        keys[i] = k
        i++
    }
    return keys
}

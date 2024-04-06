package utils

import (
	"reflect"
	"fmt"
)

func ConcatenateTwoArrays(arr1 []interface{}, arr2 []interface{}) ([]interface{}, error) {
	if reflect.TypeOf(arr1) != reflect.TypeOf(arr2) {
		fmt.Printf("Failed Concatenation, first array is of type %T, second array is of type %T", arr1, arr2)
		return nil, fmt.Errorf("Failed Concatenation, first array is of type %T, second array is of type %T", arr1, arr2)
	}
	return append(arr1, arr2...), nil
}

func RemoveDuplicate[T comparable](sliceList []T) []T {
    allKeys := make(map[T]bool)
    list := []T{}
    for _, item := range sliceList {
        if _, value := allKeys[item]; !value {
            allKeys[item] = true
            list = append(list, item)
        }
    }
    return list
}
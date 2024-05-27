package utils

import (
	"fmt"
	"reflect"
	"unicode"
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

func CompareSlices(slice1, slice2 []int) bool {
	if len(slice1) != len(slice2) {
		return false
	}
	for _, v := range slice1 {
		if !intInIntArray(v, slice2) {
			return false
		}
	}
	return true
}

func intInIntArray(n int, arr []int) bool {
	for _, i := range arr {
		if n == i {
			return true
		}
	}
	return false
}

func StringInStringArray(n string, arr []string) bool {
	for _, i := range arr {
		if n == i {
			return true
		}
	}
	return false
}
func RuneInRuneArray(n rune, arr []rune) bool {
	for _, i := range arr {
		if n == i {
			return true
		}
	}
	return false
}

func AppendStringIfNotInArr(n string, arr []string) []string {
	if !StringInStringArray(n, arr) {
		return append(arr, n)
	}
	return arr
}
func AppendRuneIfNotInArr(n rune, arr []rune) []rune {
	if !RuneInRuneArray(n, arr) {
		return append(arr, n)
	}
	return arr
}

func RuneType(r rune) string{
	switch {
	case unicode.IsUpper(r):
		return "upper"
	case unicode.IsLower(r):
		return "lower"
	case unicode.IsNumber(r):
		return "number"
	default:
		return "other"
	}
}


func DeepCopyList(original *DoublyLinkedList) *DoublyLinkedList {
    if original.Head == nil {
        return &DoublyLinkedList{}
    }

    copiedList := &DoublyLinkedList{}
    originalNode := original.Head
    var lastCopiedNode *LinkedNode = nil

    // Copia cada nodo de la lista original a la nueva lista
    for originalNode != nil {
        newNode := &LinkedNode{Value: originalNode.Value}
        if lastCopiedNode == nil {
            copiedList.Head = newNode
        } else {
            lastCopiedNode.Next = newNode
            newNode.Prev = lastCopiedNode
        }
        lastCopiedNode = newNode
        originalNode = originalNode.Next
    }

    copiedList.Tail = lastCopiedNode
    return copiedList
}

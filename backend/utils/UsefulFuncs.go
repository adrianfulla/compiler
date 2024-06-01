package utils

import (
	"fmt"
	"reflect"
	"unicode"
	"strings"
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


// StateExists verifica si un estado LRState ya existe en una lista de estados
func StateExists(states []*LRState, newState *LRState) bool {
    for _, state := range states {
        if ItemsEqual(state.Items, newState.Items) {
            return true
        }
    }
    return false
}

// ItemsEqual compara dos slices de *Item para determinar si son iguales
func ItemsEqual(items1, items2 []*Item) bool {
    if len(items1) != len(items2) {
        return false
    }
    itemMap := make(map[string]int)

    for _, item := range items1 {
        key := itemKey(item)
        itemMap[key]++
    }

    for _, item := range items2 {
        key := itemKey(item)
        if count, found := itemMap[key]; !found || count == 0 {
            return false
        } else {
            itemMap[key]--
        }
    }

    return true
}

// itemKey genera una clave única para un Item basada en su producción y posiciones del punto
func itemKey(item *Item) string {
    // Esto asume que Production tiene un identificador único o se puede representar de manera única con Head.
    // Ajusta esta implementación según la estructura de tus datos.
    return fmt.Sprintf("%s-%d-%d", item.Production.Head, item.Position, item.SubPos)
}


// IsTerminal verifica si un símbolo es terminal basado en tu convención de mayúsculas.
func IsTerminal(symbol string) bool {
    return strings.ToUpper(symbol) == symbol
}

// Contains verifica si un slice contiene un cierto valor.
func Contains(slice []string, value string) bool {
    for _, item := range slice {
        if item == value {
            return true
        }
    }
    return false
}

// Unique elimina duplicados de un slice de strings.
func Unique(slice []string) []string {
    keys := make(map[string]bool)
    list := []string{}
    for _, entry := range slice {
        if _, value := keys[entry]; !value {
            keys[entry] = true
            list = append(list, entry)
        }
    }
    return list
}

// Filter devuelve un slice filtrado según la función condicional.
func Filter(slice []string, condition func(string) bool) []string {
    var result []string
    for _, item := range slice {
        if condition(item) {
            result = append(result, item)
        }
    }
    return result
}


func BoolsToBytes(t []bool) []byte {
    b := make([]byte, (len(t)+7)/8)
    for i, x := range t {
        if x {
            b[i/8] |= 0x80 >> uint(i%8)
        }
    }
    return b
}
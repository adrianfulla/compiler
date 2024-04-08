package automatas

import (
	"encoding/json"
	"unicode"

	// "strings"
	"fmt"

	"github.com/adrianfulla/compiler/backend/utils"
)

func NewFullNodo(Valor string, Izquierdo *utils.Nodo, Derecho *utils.Nodo, Nullability bool, firstpos []int, lastpos []int) *utils.Nodo {
	return &utils.Nodo{
		Valor:       Valor,
		Izquierdo:   Izquierdo,
		Derecho:     Derecho,
		Leaf:        nil,
		Nullability: Nullability,
		Firstpos:    firstpos,
		Lastpos:     lastpos,
		Followpos:   make([]int, 0),
	}
}
func NewStarNodo(Valor string, Izquierdo *utils.Nodo, Nullability bool, firstpos []int, lastpos []int) *utils.Nodo {
	return &utils.Nodo{
		Valor:       Valor,
		Izquierdo:   Izquierdo,
		Derecho:     nil,
		Leaf:        nil,
		Nullability: Nullability,
		Firstpos:    firstpos,
		Lastpos:     lastpos,
		Followpos:   make([]int, 0),
	}
}

// alphanum determina si un carácter es alfanumérico, un epsilon o un #
func alphanum(a rune) bool {
	return unicode.IsLetter(a) || unicode.IsDigit(a) || a == 'ε' || a == '#'
}

type ArbolExpresion struct {
	Raiz     *utils.Nodo   `json:"raiz"`
	Simbolos []*utils.Nodo `json:"simbolos"`
}

func (arbol *ArbolExpresion) ConstruirArbol(posfix string) {
	stack := utils.NewStack()

	for _, char := range posfix {
		canBeNull := false
		firstpos := []int{}
		lastpos := []int{}

		if alphanum(char) {
			if char == 'ε' {
				canBeNull = true
			} else if unicode.IsLetter(char) || unicode.IsDigit(char) || char == '#' {
				firstpos = append(firstpos, len(arbol.Simbolos))
				lastpos = append(lastpos, len(arbol.Simbolos))
				stack.Push(arbol.createLeaf(string(char), canBeNull, firstpos, lastpos))
			}
		} else if char == '*' {
			canBeNull = true
			n1 := stack.Pop().(*utils.Nodo)
			stack.Push(NewStarNodo(string(char), n1, canBeNull, n1.Firstpos, n1.Lastpos))

		} else if char == '|' || char == '^' {
			n2 := stack.Pop().(*utils.Nodo)
			n1 := stack.Pop().(*utils.Nodo)
			if char == '|' {
				firstpos := append(n1.Firstpos, n2.Firstpos...)
				lastpos := append(n1.Lastpos, n2.Lastpos...)
				canBeNull := n1.Nullability || n2.Nullability
				stack.Push(NewFullNodo(string(char), n1, n2, canBeNull, firstpos, lastpos))
			} else {
				firstpos := n1.Firstpos
				lastpos := n2.Lastpos
				canBeNull := n1.Nullability && n2.Nullability
				if n1.Nullability {
					firstpos = append(firstpos, n2.Firstpos...)
				}
				if n2.Nullability {
					lastpos = append(lastpos, n1.Lastpos...)
				}
				stack.Push(NewFullNodo(string(char), n1, n2, canBeNull, firstpos, lastpos))
			}
		}
	}
	arbol.Raiz = stack.Pop().(*utils.Nodo)
	arbol.calcular_followpos()
	// arbol.imprimirDetalle()
}

func (arbol *ArbolExpresion) createLeaf(valor string, nullable bool, firstpos []int, lastpos []int) *utils.Nodo {
	if lastpos == nil {
		lastpos = []int{}
	}
	if firstpos == nil {
		firstpos = []int{}
	}
	leafid := len(arbol.Simbolos)
	nodo := &utils.Nodo{
		Valor:       valor,
		Izquierdo:   nil,
		Derecho:     nil,
		Leaf:        &leafid,
		Nullability: nullable,
		Firstpos:    firstpos,
		Lastpos:     lastpos,
		Followpos:   []int{},
	}
	arbol.Simbolos = append(arbol.Simbolos, nodo)
	return nodo
}

func (arbol *ArbolExpresion) ToJson() ([]byte, error) {
	return json.MarshalIndent(arbol, "", "")
}

func (arbol *ArbolExpresion) visitNodo(nodo *utils.Nodo) {
	if nodo.Valor == "^" {
		for _, pos := range nodo.Izquierdo.Lastpos {
			if arbol.Simbolos[pos].Followpos == nil {
				arbol.Simbolos[pos].Followpos = make([]int, 0)
			}

			arbol.Simbolos[pos].Followpos = append(arbol.Simbolos[pos].Followpos, nodo.Derecho.Firstpos...)
		}
	} else if nodo.Valor == "*" {
		for _, pos := range nodo.Lastpos {
			if arbol.Simbolos[pos].Followpos == nil {
				arbol.Simbolos[pos].Followpos = make([]int, 0)
			}

			arbol.Simbolos[pos].Followpos = append(arbol.Simbolos[pos].Followpos, nodo.Firstpos...)
		}
	}
	if nodo.Izquierdo != nil {
		arbol.visitNodo(nodo.Izquierdo)
	}
	if nodo.Derecho != nil {
		arbol.visitNodo(nodo.Derecho)
	}
}

func (arbol *ArbolExpresion) calcular_followpos() {
	for _, simbolo := range arbol.Simbolos {
		simbolo.Followpos = make([]int, 0)
	}
	arbol.visitNodo(arbol.Raiz)
}

func (arbol *ArbolExpresion) imprimirDetalle() {
	fmt.Println("Mostrando detalles de arbol")
	arbol.Raiz.ImprimirDetalle()
}

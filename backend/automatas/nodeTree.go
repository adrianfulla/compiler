package automatas

import (
	"encoding/json"
	"unicode"

	// "strings"
	"fmt"

	"github.com/adrianfulla/compiler/backend/utils"
)

func NewFullNodo(Valor rune, Izquierdo *utils.Nodo, Derecho *utils.Nodo, Nullability bool, firstpos []int, lastpos []int) *utils.Nodo {
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
func NewStarNodo(Valor rune, Izquierdo *utils.Nodo, Nullability bool, firstpos []int, lastpos []int) *utils.Nodo {
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
				stack.Push(arbol.createLeaf(char, canBeNull, firstpos, lastpos))
			}
		} else if char == '*' {
			canBeNull = true
			n1 := stack.Pop().(*utils.Nodo)
			stack.Push(NewStarNodo(char, n1, canBeNull, n1.Firstpos, n1.Lastpos))

		} else if char == '|' || char == '^' {
			n2 := stack.Pop().(*utils.Nodo)
			n1 := stack.Pop().(*utils.Nodo)
			if char == '|' {
				firstpos := append(n1.Firstpos, n2.Firstpos...)
				lastpos := append(n1.Lastpos, n2.Lastpos...)
				canBeNull := n1.Nullability || n2.Nullability
				stack.Push(NewFullNodo(char, n1, n2, canBeNull, firstpos, lastpos))
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
				stack.Push(NewFullNodo(char, n1, n2, canBeNull, firstpos, lastpos))
			}
		}
	}
	arbol.Raiz = stack.Pop().(*utils.Nodo)
	arbol.calcular_followpos()
	// arbol.imprimirDetalle()
}

func (arbol *ArbolExpresion) createLeaf(valor rune, nullable bool, firstpos []int, lastpos []int) *utils.Nodo {
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
	if nodo.Valor == '^' && nodo.Leaf == nil {
		// fmt.Printf("Nodo: %s,first: %d, last: %d \n",nodo.Valor, nodo.Firstpos, nodo.Lastpos)
		for _, pos := range nodo.Izquierdo.Lastpos {
			if arbol.Simbolos[pos].Followpos == nil {
				arbol.Simbolos[pos].Followpos = make([]int, 0)
			}

			arbol.Simbolos[pos].Followpos = append(arbol.Simbolos[pos].Followpos, nodo.Derecho.Firstpos...)
		}
	} else if nodo.Valor == '*' && nodo.Leaf == nil  {
		// fmt.Printf("Nodo: %s,first: %d, last: %d \n",nodo.Valor, nodo.Firstpos, nodo.Lastpos)
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


func (arbol *ArbolExpresion) ExtendedConstruirArbol(posfix []utils.RegexToken) {
	stack := utils.NewStack()
	// fmt.Print("ACA\n")
	for _, char := range posfix{
		canBeNull := false
		firstpos := []int{}
		lastpos := []int{}
		// fmt.Print(char)
		if char.IsOperator == "" || char.IsOperator == "NULL" || char.IsOperator == "ENDOFTREE"{
			if char.IsOperator == "NULL"{
				canBeNull = true
			}else{
				firstpos = append(firstpos, len(arbol.Simbolos))
				lastpos = append(lastpos, len(arbol.Simbolos))
				stack.Push(arbol.createLeaf(char.Value[0], canBeNull, firstpos, lastpos))
			}
		}else if char.IsOperator == "KLEENE"{
			canBeNull = true
			n1 := stack.Pop().(*utils.Nodo)
			stack.Push(NewStarNodo(char.Value[0],n1, canBeNull, n1.Firstpos, n1.Lastpos))
		}else if char.IsOperator == "OROPERATOR"|| char.IsOperator == "CATOPERATOR"{
			n2 := &utils.Nodo{}
			if stack.Peek() != nil{
				n2 = stack.Pop().(*utils.Nodo)
			}
			// else{
				
			// }
			n1 := &utils.Nodo{}
			if stack.Peek() != nil{
				n1 = stack.Pop().(*utils.Nodo)
			}else{
				n2.Nullability = true
				stack.Push(n2)
				continue
			}
			if char.IsOperator == "OROPERATOR" {
				firstpos := append(n1.Firstpos, n2.Firstpos...)
				lastpos := append(n1.Lastpos, n2.Lastpos...)
				canBeNull := n1.Nullability || n2.Nullability
				stack.Push(NewFullNodo(char.Value[0], n1, n2, canBeNull, firstpos, lastpos))
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
				stack.Push(NewFullNodo(char.Value[0], n1, n2, canBeNull, firstpos, lastpos))
			}
		}
	}
	arbol.Raiz = stack.Pop().(*utils.Nodo)
	
	arbol.calcular_followpos()
	// arbol.imprimirDetalle()
}

package automatas

import (
	"encoding/json"
	"unicode"
	// "strings"
	"fmt"
)

type Nodo struct {
	Valor       rune                   `json:"valor"`
	Izquierdo   *Nodo                    `json:"izquierdo,omitempty"`
	Derecho     *Nodo                    `json:"derecho,omitempty"`
	Leaf        int                      `json:"leaf"`
	Nullability bool                     `json:"nullability,omitempty"`
	Firstpos    []int			         `json:"firstpos,omitempty"`
	Lastpos     []int	      		     `json:"lastpos,omitempty"`
	Followpos   []int	 				 `json:"followpos,omitempty"`
}


// alphanum determina si un carácter es alfanumérico, un epsilon o un #
func alphanum(a rune) bool {
	return unicode.IsLetter(a) || unicode.IsDigit(a) || a == 'ε' || a == '#'
}

type ArbolExpresion struct {
	Raiz     *Nodo `json:"raiz"`
	Simbolos []*Nodo `json:"simbolos"`
}

func (arbol *ArbolExpresion) ConstruirArbol(posfix string) {
	stack := []*Nodo{}

	for _, char := range posfix {
		fmt.Print(char)
		canBeNull := false
		firstpos := []int{}
		lastpos := []int{}

		if alphanum(char) {
			if char == 'ε'{
				canBeNull = true
			} else if unicode.IsLetter(char) || unicode.IsDigit(char) || char == '#' {
				firstpos = append(firstpos, len(arbol.Simbolos))
                lastpos = append(lastpos, len(arbol.Simbolos))
				stack = append(stack, arbol.createLeaf(char, canBeNull, firstpos, lastpos))
			}
		} else if char == '*' {
			canBeNull = true
			n1 := stack.Pop

		} else if char == '|' || char == '^'{

		}
	}

}

func (arbol *ArbolExpresion) createLeaf(valor rune, nullable bool, firstpos []int, lastpos []int) *Nodo {
	if lastpos == nil {lastpos = []int{}}
	if firstpos == nil {firstpos = []int{}}
	nodo := &Nodo{
		Valor:       valor,
        Izquierdo:   nil,
        Derecho:     nil,
        Leaf:        len(arbol.Simbolos),
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

// func RecorrerArbol(nodo *Nodo, nivel int) {
// 	if nodo == nil {
// 		return
// 	}

// 	// Imprimir los detalles del nodo actual
// 	indent := strings.Repeat("  ", nivel) // Indentación basada en el nivel del árbol
// 	fmt.Printf("%sNodo: Valor=%s, Leaf=%d, Nullability=%t\n", indent, nodo.Valor, nodo.Leaf, nodo.Nullability)
// 	fmt.Printf("%sFirstpos=%v, Lastpos=%v, Followpos=%v\n", indent, nodo.Firstpos, nodo.Lastpos, nodo.Followpos)

// 	// Recursivamente visitar los nodos izquierdo y derecho
// 	fmt.Printf("%sIzquierdo:\n", indent)
// 	RecorrerArbol(nodo.Izquierdo, nivel+1)
// 	fmt.Printf("%sDerecho:\n", indent)
// 	RecorrerArbol(nodo.Derecho, nivel+1)
// }

// ImprimirArbol inicia el recorrido del árbol desde la raíz
func (arbol *ArbolExpresion) ImprimirArbol() {
	fmt.Println("Recorriendo el Árbol de Expresión:")
	RecorrerArbol(arbol.Raiz, 0)
}

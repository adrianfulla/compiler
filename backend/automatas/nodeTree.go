package automatas

import (
	"encoding/json"
	"unicode"
)

type Nodo struct {
	Valor       string                   `json:"valor"`
	Izquierdo   *Nodo                    `json:"izquierdo,omitempty"`
	Derecho     *Nodo                    `json:"derecho,omitempty"`
	Leaf        int                      `json:"leaf"`
	Nullability bool                     `json:"nullability,omitempty"`
	Firstpos    map[int]struct{}         `json:"firstpos,omitempty"`
	Lastpos     map[int]struct{}         `json:"lastpos,omitempty"`
	Followpos   map[int]map[int]struct{} `json:"followpos,omitempty"`
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
		if alphanum(char) || char == 'ε' {
			// Si es un operando, crea un nodo y lo apila
			nodo := &Nodo{
				Valor:    string(char),
				Firstpos: make(map[int]struct{}),
				Lastpos:  make(map[int]struct{}),
			}
			stack = append(stack, nodo)
		} else {
			// Es un operador, pop elementos de la pila y crea nuevos nodos
			switch char {
			case '*':
				nodo := &Nodo{
					Valor: string(char),
				}
				if len(stack) > 0 {
					nodo.Izquierdo = stack[len(stack)-1]
					stack = stack[:len(stack)-1]
				}
				stack = append(stack, nodo)
			case '|', '^':
				nodo := &Nodo{
					Valor: string(char),
				}
				if len(stack) > 1 {
					nodo.Derecho = stack[len(stack)-1]
					nodo.Izquierdo = stack[len(stack)-2]
					stack = stack[:len(stack)-2]
				}
				stack = append(stack, nodo)
				// Agregar más casos según sea necesario
			}
		}
	}

	if len(stack) > 0 {
		arbol.Raiz = stack[0]
	}
}

func (arbol *ArbolExpresion) ToJson() ([]byte, error) {
	return json.MarshalIndent(arbol, "", "")
}

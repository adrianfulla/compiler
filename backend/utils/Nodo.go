package utils

import (
	"fmt"
)

type Nodo struct {
	Valor       string `json:"valor"`
	Izquierdo   *Nodo  `json:"izquierdo,omitempty"`
	Derecho     *Nodo  `json:"derecho,omitempty"`
	Leaf        *int   `json:"leaf"`
	Nullability bool   `json:"nullability,omitempty"`
	Firstpos    []int  `json:"firstpos,omitempty"`
	Lastpos     []int  `json:"lastpos,omitempty"`
	Followpos   []int  `json:"followpos,omitempty"`
}

func (nodo *Nodo) IsLeaf() bool {
	return nodo.Leaf != nil
}

func (nodo *Nodo) PrintNodoDetalle() {
	if nodo.IsLeaf() {
		fmt.Printf("Nodo %d\n", *nodo.Leaf)
	}
	fmt.Printf("Valor: %s\n", string(nodo.Valor))
	fmt.Printf("Nullability: %t\n", nodo.Nullability)
	fmt.Printf("Firstpos: %d\n", nodo.Firstpos)
	fmt.Printf("Lastpos: %d\n", nodo.Lastpos)
	fmt.Printf("Followpos: %d\n", nodo.Followpos)
}

func (nodo *Nodo) ImprimirDetalle() {
	nodo.PrintNodoDetalle()
	if !nodo.IsLeaf() {
		if nodo.Derecho != nil {
			fmt.Println("Nodo Derecho")
			nodo.Derecho.ImprimirDetalle()
		}
		if nodo.Izquierdo != nil {
			fmt.Println("Nodo Izquierdo")
			nodo.Izquierdo.ImprimirDetalle()
		}
	}
}

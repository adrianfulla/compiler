package automatas

import (
	"fmt"
	"strconv"

	"strings"

	"github.com/adrianfulla/compiler/backend/utils"
)

type Dstate struct {
	nombre       string
	aceptacion   bool
	transiciones map[rune]*Dstate
	posicion     int
}

func NewDstate(nombre string, aceptacion bool, posicion int) *Dstate {
	return &Dstate{
		nombre:       nombre,
		aceptacion:   aceptacion,
		transiciones: make(map[rune]*Dstate),
		posicion:     posicion,
	}
}

func (d *Dstate) AddTransicion(simbolo rune, estado *Dstate) {
	d.transiciones[simbolo] = estado
}

type DirectAfd struct {
	transiciones      map[string]map[string]string
	estadosAceptacion map[string]string
	alfabeto          map[string]string
	estados           map[string]*Dstate
	estadoInicial     *Dstate
	arbol             *ArbolExpresion
	estadoActual      int
	posiciones        []int
}

func intInIntArray(n interface{}, arr []int) bool {
    switch v := n.(type) {
    case int:
        for _, i := range arr {
            if i == v {
                return true
            }
        }
    case *int:
        if v == nil {
            return false
        }
        for _, i := range arr {
            if i == *v {
                return true
            }
        }
    default:
        // Si n no es ni int ni *int, se retorna false
        return false
    }
    return false
}

func simboloInSimboloDict(n rune, dict map[rune][]int)(bool){
	for i := range dict{
		if n == i{
			return true
		}
	}
	return false
}

func removeDuplicates(arr []interface{}){

}

// func (afd *DirectAfd) imprimirDetalle() {
// 	fmt.Println("Detalle del AFD:")
// 	fmt.Println("Transiciones:")
// 	for estado, transMap := range afd.transiciones {
// 		for simbolo, destino := range transMap {
// 			fmt.Printf("  %s -> %s: %s\n", estado, simbolo, destino)
// 		}
// 	}

// 	fmt.Println("Estados de Aceptación:")
// 	for estado, aceptacion := range afd.estadosAceptacion {
// 		fmt.Printf("  %s: %t\n", estado, aceptacion)
// 	}

// 	fmt.Println("Alfabeto:")
// 	for simbolo := range afd.alfabeto {
// 		fmt.Printf("  %s\n", simbolo)
// 	}

// 	fmt.Println("Estados:")
// 	for nombre, estado := range afd.estados {
// 		fmt.Printf("  %s: %v\n", nombre, estado)
// 	}

// 	if afd.estadoInicial != nil {
// 		fmt.Printf("Estado Inicial: %s\n", afd.estadoInicial.nombre)
// 	} else {
// 		fmt.Println("Estado Inicial: No definido")
// 	}
// }

func NewDirectAfd(regex string) *DirectAfd {
	afd := &DirectAfd{
		transiciones:      make(map[string]map[string]string),
		estadosAceptacion: make(map[string]string),
		alfabeto:          make(map[string]string),
		estados:           make(map[string]*Dstate),
	}
	afd.estadoActual = 0
	afd.arbol = &ArbolExpresion{}
	afd.arbol.ConstruirArbol(regex + "#^")
	afd.construirAfd()

	// afd.imprimirDetalle()
	return afd
}

func (afd *DirectAfd) nuevoEstado(position int, aceptacion bool) *Dstate {
	nombre := "S" + strconv.Itoa(afd.estadoActual)
	afd.estadoActual++
	nuevo_estado := NewDstate(nombre, aceptacion, position)
	afd.estados[nombre] = nuevo_estado
	afd.posiciones = append(afd.posiciones, position)
	return nuevo_estado
}

func (afd *DirectAfd) obtenerOCrearEstado(positions []int) *Dstate {
	for _, estado := range afd.estados {
		if intInIntArray(estado.posicion, positions) {
			return estado
		}
	}
	aceptacion := intInIntArray(afd.arbol.Raiz.Derecho.Leaf,positions)
	return afd.nuevoEstado(len(positions), aceptacion)
}

// construirAfd construye el AFD a partir del árbol de expresión.
func (afd *DirectAfd) construirAfd() {
	fmt.Println()
	afd.estadoInicial = afd.obtenerOCrearEstado(afd.arbol.Raiz.Firstpos)
	pendientes := utils.NewStack()
	pendientes.Push(afd.estadoInicial)
	procesados := utils.NewStack()

	for pendientes.Size() > 0 {
		curr_estado := pendientes.Pop().(*Dstate)
		fmt.Printf("curr est: %s \n", curr_estado.nombre)
		if !procesados.ElemInStack(curr_estado.nombre) {
			simbolos_a_pos := make(map[rune][]int)
			for pos := range curr_estado.posicion {
				fmt.Printf("pos %d ", pos)
				simbolo := afd.arbol.Simbolos[pos].Valor
				fmt.Printf("simbolo %s ", string(simbolo))
				if !strings.ContainsRune("ε#", simbolo) {
					followpos := afd.arbol.Simbolos[pos].Followpos
					fmt.Printf("followpos %d\n", followpos)
						if simboloInSimboloDict(simbolo, simbolos_a_pos){
							simbolos_a_pos[simbolo] = utils.RemoveDuplicate(append(simbolos_a_pos[simbolo], followpos...))
						}else{
							simbolos_a_pos[simbolo] = followpos
						}
				}
			}
			for sim, pos := range simbolos_a_pos{
				next_state := afd.obtenerOCrearEstado(pos)
				curr_estado.AddTransicion(sim, next_state)
				if !procesados.ElemInStack(next_state.nombre) && !pendientes.ElemInStack(next_state){
					pendientes.Push(next_state)
				}
			}
		}
	}
}
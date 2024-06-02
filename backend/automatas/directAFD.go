package automatas

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/adrianfulla/compiler/backend/utils"
	// "time"
)

type Dstate struct {
	nombre       string
	aceptacion   bool
	transiciones map[rune]*Dstate
	posicion     []int
}

func NewDstate(nombre string, aceptacion bool, posicion []int) *Dstate {
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

func (d *Dstate) Print() {
	fmt.Printf("Dstate: [nombre: %s, aceptacion: %t, posicion: %d, transiciones: [", d.nombre, d.aceptacion, d.posicion)
	for i, e := range d.transiciones {
		fmt.Printf("%s: %s, ", string(i), e.nombre)
	}
	fmt.Print("]]\n")
}

func printDstateStack(stack utils.Stack) {
	temp := stack
	for temp.Size() > 0 {
		t := temp.Pop()
		switch t.(type) {
		case *Dstate:
			t.(*Dstate).Print()
		case string:
			fmt.Printf("String found %s\n", t)
		default:
			fmt.Printf("Dstate not in stack, obtained elem of type %t\n", t)
		}
	}
}

type DirectAfd struct {
	transiciones      map[string]map[string]string
	estadosAceptacion []string
	alfabeto          []string
	estados           map[string]*Dstate
	estadoInicial     *Dstate
	Arbol             *ArbolExpresion
	estadoActual      int
	posiciones        map[int][]int
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

func simboloInSimboloDict(n rune, dict map[rune][]int) bool {
	for i := range dict {
		if n == i {
			return true
		}
	}
	return false
}

func NewDirectAfd(regex string) *DirectAfd {
	afd := &DirectAfd{
		transiciones:      make(map[string]map[string]string),
		estadosAceptacion: []string{},
		alfabeto:          []string{},
		estados:           make(map[string]*Dstate),
		posiciones:        make(map[int][]int),
	}
	afd.estadoActual = 0
	afd.Arbol = &ArbolExpresion{}
	afd.Arbol.ConstruirArbol(regex + "#^")
	// afd.Arbol.imprimirDetalle()
	afd.construirAfd()

	return afd
}

func ExtendedNewDirectAfd(regex []utils.RegexToken) *DirectAfd{
	afd := &DirectAfd{
		transiciones:      make(map[string]map[string]string),
		estadosAceptacion: []string{},
		alfabeto:          []string{},
		estados:           make(map[string]*Dstate),
		posiciones:        make(map[int][]int),
	}

	// for _, token := range regex{
	// 	fmt.Printf(" %s, operator: %s", token.Value, token.IsOperator)
	// }

	afd.estadoActual = 0
	afd.Arbol = &ArbolExpresion{}
	afd.Arbol.ExtendedConstruirArbol(append(regex, utils.RegexToken{
		Value: []rune{'#'},
		IsOperator: "ENDOFTREE",
	}, utils.RegexToken{
		Value: []rune{'^'},
		IsOperator: "CATOPERATOR",
	}, ))
	// afd.Arbol.imprimirDetalle()
	afd.construirAfd()
	// fmt.Print("Afd contruido\n")
	return afd
}

func (afd *DirectAfd) nuevoEstado(position []int, aceptacion bool) *Dstate {
	nombre := "S" + strconv.Itoa(afd.estadoActual)
	afd.estadoActual++
	nuevo_estado := NewDstate(nombre, aceptacion, position)
	afd.estados[nombre] = nuevo_estado
	afd.posiciones[len(afd.posiciones)] = position
	return nuevo_estado
}

func (afd *DirectAfd) obtenerOCrearEstado(positions []int) *Dstate {
    for _, estado := range afd.estados {
        if utils.CompareSlices(estado.posicion, positions) {
            return estado
        }
    }
    // Identificar si el estado debe ser de aceptación
    aceptacion := false
    // Un estado es de aceptación si alguna de las posiciones puede ser el final de la expresión regular
    for _, pos := range positions {
        if afd.Arbol.Simbolos[pos].Valor == '#' { // Suponiendo que '#' es el carácter de finalización
            aceptacion = true
            break
        }
    }

    return afd.nuevoEstado(positions, aceptacion)
}



func (afd *DirectAfd) construirAfd() {
    afd.estadoInicial = afd.obtenerOCrearEstado(afd.Arbol.Raiz.Firstpos)
    pendientes := utils.NewStack()
    pendientes.Push(afd.estadoInicial)
    procesados := make(map[string]bool)

    for pendientes.Size() > 0 {
		// fmt.Printf("Pendientes: %d\n", pendientes.Size())
        curr_estado := pendientes.Pop().(*Dstate)
        if !procesados[curr_estado.nombre] {
            simbolos_a_pos := make(map[rune][]int)
            for _, pos := range curr_estado.posicion {
                simbolo := afd.Arbol.Simbolos[pos].Valor
                if !strings.ContainsRune("ε#", simbolo) {
                    simbolos_a_pos[simbolo] = append(simbolos_a_pos[simbolo], afd.Arbol.Simbolos[pos].Followpos...)
                }
            }

            for sim, pos := range simbolos_a_pos {
				next_state := afd.obtenerOCrearEstado(pos)
				curr_estado.AddTransicion(sim, next_state)
				if !procesados[next_state.nombre] && !pendientes.ElemInStack(next_state) {
					pendientes.Push(next_state)
				}
			}

            procesados[curr_estado.nombre] = true
        }
    }
}

func (afd *DirectAfd) MarshalJson() ([]byte, error) {
	return json.Marshal(afd.ToJson())
}

func (afd *DirectAfd) ToJson() *DAfdJson {
	jsonMaker := &DAfdJson{
		Estados:        []string{},
		Alfabeto:       []rune{},
		EstadosFinales: []string{},
		Transiciones:   make(map[string]map[string]string),
	}
	jsonMaker.EstadoInicial = afd.estadoInicial.nombre
	for _, estado := range afd.estados {
		jsonMaker.Estados = utils.AppendStringIfNotInArr(estado.nombre, jsonMaker.Estados)
		if estado.aceptacion {
			jsonMaker.EstadosFinales = utils.AppendStringIfNotInArr(estado.nombre, jsonMaker.EstadosFinales)
		}
		for sim, trans := range estado.transiciones {
			jsonMaker.Alfabeto = utils.AppendRuneIfNotInArr(sim, jsonMaker.Alfabeto)
			if jsonMaker.Transiciones[estado.nombre] == nil {
				jsonMaker.Transiciones[estado.nombre] = make(map[string]string)
			}
			jsonMaker.Transiciones[estado.nombre][string(sim)] = trans.nombre
		}
	}
	return jsonMaker
}

type DAfdJson struct {
	Estados        []string                     `json:"estados"`
	Alfabeto       []rune                     	`json:"alfabeto"`
	EstadoInicial  string                       `json:"estado_inicial"`
	EstadosFinales []string                     `json:"estados_finales"`
	Transiciones   map[string]map[string]string `json:"transiciones"`
}

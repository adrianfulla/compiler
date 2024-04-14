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
	transiciones map[string]*Dstate
	posicion     []int
}

func NewDstate(nombre string, aceptacion bool, posicion []int) *Dstate {
	return &Dstate{
		nombre:       nombre,
		aceptacion:   aceptacion,
		transiciones: make(map[string]*Dstate),
		posicion:     posicion,
	}
}

func (d *Dstate) AddTransicion(simbolo string, estado *Dstate) {
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

func ExtendedNewDirectAfdI(regex []utils.RegexToken) *DirectAfd{
	afd := &DirectAfd{
		transiciones:      make(map[string]map[string]string),
		estadosAceptacion: []string{},
		alfabeto:          []string{},
		estados:           make(map[string]*Dstate),
		posiciones:        make(map[int][]int),
	}
	afd.estadoActual = 0
	afd.Arbol = &ArbolExpresion{}
	afd.Arbol.ExtendedConstruirArbol(append(regex, utils.RegexToken{
		Value: []rune{'#'},
		IsOperator: "ENDOFTREE",
	}))
	// afd.Arbol.imprimirDetalle()
	afd.construirAfd()

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
	aceptacion := intInIntArray(*afd.Arbol.Raiz.Derecho.Leaf, positions)
	return afd.nuevoEstado(positions, aceptacion)
}

// construirAfd construye el AFD a partir del árbol de expresión.
func (afd *DirectAfd) construirAfd() {
	afd.estadoInicial = afd.obtenerOCrearEstado(afd.Arbol.Raiz.Firstpos)
	pendientes := utils.NewStack()
	pendientes.Push(afd.estadoInicial)
	procesados := utils.NewStack()

	// fmt.Print("Simbolos del arbol\n")
	// for x, y := range afd.Arbol.Simbolos{
	// 	fmt.Printf("Posicion %d con valor %s, followpos %d\n", x, y.Valor, y.Followpos)
	// }
	for pendientes.Size() > 0 {
		curr_estado := pendientes.Pop().(*Dstate)
		if !procesados.ElemInStack(curr_estado.nombre) {
			simbolos_a_pos := make(map[string][]int)
			for _, pos := range curr_estado.posicion {
				// fmt.Printf("indice %d Posicion %d\n", xy, pos)
				simbolo := afd.Arbol.Simbolos[pos].Valor
				if !strings.ContainsAny("ε#", simbolo) {
					if simbolos_a_pos[afd.Arbol.Simbolos[pos].Valor] == nil {
						simbolos_a_pos[afd.Arbol.Simbolos[pos].Valor] = make([]int, 0)
					}
					simbolos_a_pos[afd.Arbol.Simbolos[pos].Valor] = append(simbolos_a_pos[afd.Arbol.Simbolos[pos].Valor], afd.Arbol.Simbolos[pos].Followpos...)
					// fmt.Printf("Simbolos a pos: [simbolo:%s, pos: %d, followpos: %d]\n", simbolo, simbolos_a_pos[afd.Arbol.Simbolos[pos].Valor], afd.Arbol.Simbolos[pos].Followpos)
				}
			}

			for sim, pos := range simbolos_a_pos {
				next_state := afd.obtenerOCrearEstado(pos)
				curr_estado.AddTransicion(sim, next_state)
				if !procesados.ElemInStack(next_state.nombre) && !pendientes.ElemInStack(next_state) {
					pendientes.Push(next_state)
				}
			}
			procesados.Push(curr_estado.nombre)
			// fmt.Print("\nPendientes\n")
			// printDstateStack(*pendientes)
			// fmt.Print("\nProcesados\n")
			// printDstateStack(*procesados)
			// fmt.Print("\nEstaods\n")
			// for val, state := range afd.estados {
			// 	fmt.Printf("Estado %s con nombre %s tiene %d con transiciones\n", val, state.nombre, state.posicion)
			// 	for sim, trans := range state.transiciones {
			// 		fmt.Printf("Simbolo %s tiene transicion a %s\n", sim, trans.nombre)
			// 	}
			// }
		}
		// time.Sleep(1 * time.Second)
	}
}

func (afd *DirectAfd) MarshalJson() ([]byte, error) {
	return json.Marshal(afd.ToJson())
}

func (afd *DirectAfd) ToJson() *DAfdJson {
	jsonMaker := &DAfdJson{
		Estados:        []string{},
		Alfabeto:       []string{},
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
			jsonMaker.Alfabeto = utils.AppendStringIfNotInArr(sim, jsonMaker.Alfabeto)
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
	Alfabeto       []string                     `json:"alfabeto"`
	EstadoInicial  string                       `json:"estado_inicial"`
	EstadosFinales []string                     `json:"estados_finales"`
	Transiciones   map[string]map[string]string `json:"transiciones"`
}

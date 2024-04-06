package automatas

import(
	"fmt"
	"strconv"
)

type Dstate struct {
	nombre       string
	aceptacion   bool
	transiciones map[string]*Dstate
	posicion	 int
}

func NewDstate(nombre string, aceptacion bool, posicion int) *Dstate {
	return &Dstate{
		nombre:       nombre,
		aceptacion:   aceptacion,
		transiciones: make(map[string]*Dstate),
		posicion: posicion,
	}
}

func (d *Dstate) AddTransicion(simbolo string, estado *Dstate) {
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
	posiciones 		  []int
}

func intInIntArray(n int, arr []int) (bool){
	for _, i := range arr {
		if i == n {
            return true
			}
		}
	return false
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
	afd.arbol.ConstruirArbol(regex+"#^") 
	afd.construirAfd() 

	// afd.imprimirDetalle()
	return afd
}

func (afd *DirectAfd) nuevoEstado(position int, aceptacion bool) (*Dstate) {
	nombre := "S" + strconv.Itoa(afd.estadoActual)
	afd.estadoActual ++
	nuevo_estado := NewDstate(nombre, aceptacion, position)
	afd.estados[nombre] = nuevo_estado
	afd.posiciones = append(afd.posiciones, position)
	return nuevo_estado
}

func (afd *DirectAfd) obtenerOCrearEstado(positions []int) (*Dstate) { 
	 for _,estado := range afd.estados {
		if intInIntArray(estado.posicion,positions){
			return estado
		} 
	 }
	//  aceptacion := intInIntArray(afd.arbol.Raiz.Derecho.Leaf,positions)
	 aceptacion := false
	 return afd.nuevoEstado(len(positions), aceptacion)
}


// construirAfd construye el AFD a partir del árbol de expresión.
func (afd *DirectAfd) construirAfd() {
	fmt.Println()
//    inicial := afd.obtenerOCrearEstado(afd.arbol.Raiz.Firstpos)
}



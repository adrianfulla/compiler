package main

import (
	"fmt"
	"net/http"

	"github.com/adrianfulla/compiler/backend/automatas"
	"github.com/adrianfulla/compiler/backend/lexer"
	"github.com/gin-gonic/gin"
)

func main() {
	// serve()
	// json := makeDirAfd("(a|b)*aab")
	// afd, err := automatas.JsonToDfa(json)
	// if err != nil {
	// 	fmt.Printf("Error al pasar a DFA")
	// } else {
	// 	fmt.Printf("AFD [esIn: %s,\n esFin: %s,\n alfabeto:%s,\n transiciones:\n", afd.EstadoInicial, afd.EstadosFinales, afd.Alfabeto)
	// 	for state, transicion := range afd.Transiciones {
	// 		fmt.Printf("%s:{\n", state)
	// 		for sim, next_state := range transicion {
	// 			fmt.Printf("%s: %s\n", sim, next_state)
	// 		}
	// 		fmt.Printf("}\n")
	// 	}
	// 	fmt.Printf("]\n")
	// 	automatas.ExtendedSimulateAfd("aaabcaabb", *afd)
	// }
	file := (`(* Lexer para Gramática No. 1 - Expresiones aritméticas simples para variables *)

	(* Introducir cualquier header aqui *)
	
	let delim = [' ''\t''\n']
	let ws = delim+
	let letter = ['A'-'Z''a'-'z']
	let digit = ['0'-'9']
	let id = letter(letter|digit)*
	
	rule tokens = 
		ws
	  | id        { return ID }               (* Cambie por una acción válida, que devuelva el token *)
	  | '+'       { return PLUS }
	  | '*'       { return TIMES }
	  | '('       { return LPAREN }
	  | ')'       { return RPAREN }
	
	(* Introducir cualquier trailer aqui *)`)
	lexFile(file)

}

func serve() {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.POST("/automata/arbol", func(c *gin.Context) {
		var request struct {
			Regex string `json:"regex"`
		}

		if err := c.BindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		response := makeArboldeNodos(request.Regex)
		c.Data(http.StatusOK, "application/json", response)
	})
	r.POST("/automata/afd", func(c *gin.Context) {
		var request struct {
			Regex string `json:"regex"`
		}

		if err := c.BindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		response := makeDirAfd(request.Regex)
		c.Data(http.StatusOK, "application/json", response)
	})

	r.Run()
}

func makeDirAfd(Regex string) []byte {
	// Prueba de la función de validación
	postfix, err := automatas.InfixToPosfix(Regex)
	if err != nil {
		fmt.Println("Error:", err)
		return nil
	} else {
		afd := automatas.NewDirectAfd(postfix)
		jsonAfd, err := afd.MarshalJson()
		if err != nil {
			fmt.Println("Error al convertir a JSON:", err)
			return nil
		}

		return jsonAfd
	}
}

func makeArboldeNodos(Regex string) []byte {
	postfix, err := automatas.InfixToPosfix(Regex)
	if err != nil {
		fmt.Println("Error:", err)
		return nil
	} else {
		arbol := &automatas.ArbolExpresion{}
		arbol.ConstruirArbol(postfix)
		jsonAfd, err := arbol.ToJson()
		if err != nil {
			fmt.Println("Error al convertir a JSON:", err)
			return nil
		}

		return jsonAfd
	}
}

func lexFile(ymlFile string){
	lexer.LexYmlFile(ymlFile)
}

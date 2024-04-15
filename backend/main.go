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
	file := (`
	(* LexerparaGramaticaNo.4 *)

(* Introducircualquierheaderaqui *)

let ws = ' '+
let letter = ['A'-'Z''a'-'z']
let digit = ['0'-'9']
let digits = digit+
let id = letter+digit*

rule tokens = 
    ws
  | id        { return ID }               (* Cambieporunaaccipnvalida,quedevuelvaeltoken *)
  | ';'       { return SEMICOLON }
  | ":="      { return ASSIGNOP }
  | '<'       { return LT }
  | '='       { return EQ }
  | '+'       { return PLUS }
  | '-'       { return MINUS }
  | '*'       { return TIMES }
  | '/'       { return DIV }
  | '('       { return LPAREN }
  | ')'       { return RPAREN }

(* Introducircualquiertraileraqui *)`)

	Scanner, err := lexFile(file)
	if err != nil{
		fmt.Println(err)
	}else{
		var input string
		// fmt.Print("Input text: \n")
		// fmt.Scanln(&input)
		input = "Este"
		AcceptedExp, err := Scanner.ScanFile(input)
		if err != nil{
			fmt.Println(err)
		}
		for _, accepted := range AcceptedExp{
			fmt.Println(accepted)
		}
	}

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
	r.POST("/automata/afd/", func(c *gin.Context) {
		var request struct {
			Regex string 				`json:"regex"`
			// Afd   	`json:"afd"`
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

func lexFile(ymlFile string) (*lexer.Scanner, error){
	scanner, err := lexer.LexYmlFile(ymlFile)
	if err != nil{
		return nil, fmt.Errorf("error parsing yml file")
	}
	return scanner, nil
}

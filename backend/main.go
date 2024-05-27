package main

import (
	"fmt"
	"net/http"

	"github.com/adrianfulla/compiler/backend/automatas"
	"github.com/adrianfulla/compiler/backend/lexer"
	"github.com/adrianfulla/compiler/backend/parser"
	"github.com/gin-gonic/gin"
)

func main() {
	// serve()
	YalexFile := (`(* Yalex for reading yapars *)

	let comment = '/''*'((' '|[^'/'])*)'*''/'
	let lowercase = ['a'-'z']+
	let uppercase = 'I'['A'-'H''J'-'Z']+|['A'-'H''J'-'Z']['A'-'Z']*
	let token = '%''t''o''k''e''n'
	let ignoreflag = 'I''G''N''O''R''E'
	let twodots = ':'
	let semicolon = ';'
	let or = '|'
	let splitter = '%''%'
	let space = [' ''\t']+
	let newline = ['\n']+
	
	rule tokens = 
	  comment            { return COMMENT }
	  | lowercase        { return LOWERCASE }  
	  | uppercase        { return UPPERCASE }
	  | token            { return TOKEN }
	  | ignoreflag       { return IGNOREFLAG }
	  | twodots          { return TWODOTS }	
	  | semicolon        { return SEMICOLON}
	  | or               { return OR }
	  | splitter         { return SPLITTER }
	  | space            { return SPACE }
	  | newline          { return NEWLINE }
	
	  (* Footer *)
	`)

	YaparFile := (`/* Yapar for reading yapars */

	%token COMMENT LOWERCASE UPPERCASE TOKEN IGNOREFLAG TWODOTS SEMICOLON OR SPLITTER SPACE NEWLINE
	IGNORE SPACE
	IGNORE COMMENT
	
	%%
	
	file:
		filedeclarations SPLITTER newlines filerules
	;
	
	filedeclarations:
	  declarations
	  | newlines declarations
	;
	
	filerules:
	  rules
	  | rules newlines
	;
	
	/* Declarations section */
	declarations:
		declaration
	  | declarations declaration
	;
	
	declaration:
		tokendeclaration
	  | ignoredeclaration
	;
	
	tokendeclaration:
		TOKEN idlist newlines
	;
	
	ignoredeclaration:
		IGNOREFLAG idlist newlines
	;
	
	idlist:
		UPPERCASE
	  | idlist UPPERCASE
	;
	
	/* Rules section */
	rules:
		rule
	  | rules rulewithnewline
	;
	
	rulewithnewline:
		newlines rule
	  | rule
	;
	
	rule:
		rulename production semicoloncomposed
	;
	
	semicoloncomposed:
		SEMICOLON
	  | newlines SEMICOLON
	;
	
	rulename:
		LOWERCASE TWODOTS
	  | LOWERCASE TWODOTS newlines
	;
	
	production:
		productionterm
	  | production orcomposed productionterm
	;
	
	orcomposed:
		OR
	  | newlines OR
	;
	
	productionterm:
		idorliteral
	  | productionterm idorliteral
	;
	
	idorliteral:
		UPPERCASE
	  | LOWERCASE
	;
	
	newlines:
	  NEWLINE
	  | newlines NEWLINE
	;`)

	Scanner, err := lexFile(YalexFile)
	if err != nil{
		fmt.Println(err)
	}else{
		Scanner.PrintScanner()
		// var input string
		// fmt.Print("Input text: \n")
		// fmt.Scanln(&input)
		// input = "12 3 21 "
		// AcceptedExp, err := Scanner.ScanFile(input)
		// if err != nil{
		// 	fmt.Println(err)
		// }
		// fmt.Println("\nTokens aceptados")
		// for _, accepted := range AcceptedExp{
		// 	fmt.Println(accepted)
		
		Parser, err := lexYaparFile(YaparFile, Scanner)
		if err != nil{
			fmt.Println(err)
		}else{
			Parser.PrintParser()
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
func lexYaparFile(yaparFile string, scanner *lexer.Scanner) (*parser.Parser, error){
	parser, err := parser.LexYaparFile(yaparFile, scanner)
	if err != nil{
		return nil, fmt.Errorf("error parsing yapar file")
	}
	return parser, nil
}

package main

import (
	// "encoding/json"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/adrianfulla/compiler/backend/automatas"
	"github.com/adrianfulla/compiler/backend/lexer"
	"github.com/adrianfulla/compiler/backend/parser"
	"github.com/adrianfulla/compiler/backend/utils"
	"github.com/gin-gonic/gin"
)

func main() {
	serve() 
	YalexFile := (`(* Yalex for reading yapars *)

	let comment = '/* '['A'-'Z''a'-'z''0'-'9'" ."]*' */'
	let lowercase = ['a'-'z']+
	let uppercase = 'I'['A'-'H''J'-'Z']+|['A'-'H''J'-'Z']['A'-'Z']*
	let token = "%token"
	let ignoreflag = 'IGNORE '
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
	  | semicolon        { return SEMICOLON }
	  | or               { return OR }
	  | splitter         { return SPLITTER }
	  | space            { return SPACE }
	  | newline          { return NEWLINE }

	  (* Footer *)
	`)

	YaparFile := (`/* Yapar for reading yapars */ %token COMMENT 
	%token LOWERCASE UPPERCASE TOKEN IGNOREFLAG TWODOTS SEMICOLON OR SPLITTER SPACE NEWLINE
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
	; `)

	Scanner, err := lexFile(YalexFile)
	if err != nil{
		fmt.Println(err)
	}else{
		
	parserScanner, err := lexYaparFile(YaparFile, Scanner)
	if err != nil{
		fmt.Println(err)
	}else{
		_,table, err := makeSLR(parserScanner)
		if err != nil{
			fmt.Println(err)
		}
		val, err := parseString(parserScanner, YaparFile, table)
		if err != nil{
			fmt.Println(err)
		}
		fmt.Printf("El resultado fue %t", val)
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
		response,_, err := makeDirAfd(request.Regex)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.Data(http.StatusOK, "application/json", response)
	})
	r.POST("/automata/afd/simulate/", func(c *gin.Context) {
		var request struct {
			Regex string `json:"regex"`
			Simulate string   	`json:"simulate"`
		}
		
		if err := c.BindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		_,afd, err := makeDirAfd(request.Regex)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		response, val, err := simulateDirAfd(afd, request.Simulate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		fmt.Println(val, response)
		c.Data(http.StatusOK, "application/json", response)
	})
	r.POST("/lexer/create/", func(c *gin.Context) {
		var request struct {
			Yalex string `json:"yalex`
		}
		if err := c.BindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		response := createLexer(request.Yalex)
		c.Data(http.StatusOK, "application/json", response)
	})

	r.POST("/lexer/parselex/", func(c *gin.Context) {
		var request struct {
			Yalex   string `json:"yalex"`
			Parsing string `json:"parsing"`
		}
		if err := c.BindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		response, err := parseLex(request.Yalex, request.Parsing)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.Data(http.StatusOK, "application/json", response)
	})

	r.POST("/parser/slr/",func(c *gin.Context) {
		var request struct {
			Yalex   string `json:"yalex"`
			Yapar string `json:"yapar"`
		}
		if err := c.BindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// fmt.Print("ACA1")
		Scanner, err := lexFile(request.Yalex)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// fmt.Print("ACA2")
		parser, err := lexYaparFile(request.Yapar, Scanner)
		if err != nil {
			fmt.Print(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// fmt.Print("ACA3")
		response,_, err := makeSLR(parser)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.Data(http.StatusOK, "application/json", response)
	})

	r.POST("/parser/lr1/",func(c *gin.Context) {
		var request struct {
			Yalex   string `json:"yalex"`
			Yapar string `json:"yapar"`
			Parsing string `json:"parsing"`
		}
		if err := c.BindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// fmt.Print("ACA1")
		Scanner, err := lexFile(request.Yalex)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// fmt.Print("ACA2")
		parser, err := lexYaparFile(request.Yapar, Scanner)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// fmt.Print("ACA3")
		_,table, err := makeSLR(parser)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		response, err := parseString(parser,request.Parsing, table)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.Data(http.StatusOK, "application/json", response)
	})

	r.Run()
}

func makeDirAfd(Regex string) ([]byte, *automatas.DirectAfd, error) {
	// Prueba de la función de validación
	postfix, err := automatas.InfixToPosfix(Regex)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, nil, err
	} else {
		afd := automatas.NewDirectAfd(postfix)
		jsonAfd, err := afd.MarshalJson()
		if err != nil {
			fmt.Println("Error al convertir a JSON:", err)
			return nil, nil,err
		}

		return jsonAfd, afd,nil
	}
}

func simulateDirAfd(afd *automatas.DirectAfd, simulate string) ([]byte,string, error){
	afdJson := afd.ToJson()
	returnVal := "false"
	_,_,err := automatas.SimulateDFA(simulate, afdJson.EstadoInicial, afdJson.EstadosFinales, afdJson.Transiciones)
	if err == nil{
		returnVal = "true"
	}

	// fmt.Print(returnVal)
	response := []byte(returnVal)

	return response, returnVal, nil

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

func createLexer(ymlFile string) (result []byte) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from error:", r)
			result = []byte("false")
		}
	}()

	_, err := lexFile(ymlFile)
	if err != nil {
		return []byte("false")
	}
	return []byte("true")
}

func parseLex(ymlFile string, parseString string) ([]byte, error) {
	scanner, err := lexFile(ymlFile)
	if err != nil {
		return nil, fmt.Errorf("error parsing yml file")
	}
	accepted, err := scanner.ScanFile(parseString)
	if err != nil {
		return nil, fmt.Errorf("error parsing string")
	}
	jsonData, err := json.Marshal(accepted)
	if err != nil {
		fmt.Print("error serializing expression to  Json")
		return nil, fmt.Errorf("error serializing expression to  Json")
	}
	return jsonData, nil
}

func lexFile(ymlFile string) (*lexer.Scanner, error) {
	scanner, err := lexer.LexYmlFile(ymlFile)
	if err != nil {
		return nil, fmt.Errorf("error parsing yml file")
	}
	return scanner, nil
}
func lexYaparFile(yaparFile string, scanner *lexer.Scanner) (*parser.Parser, error) {
	parser, err := parser.LexYaparFile(yaparFile, scanner)
	if err != nil {
		return nil, err
	}
	// parser.PrintParser()
	return parser, nil
}

func makeSLR(parser *parser.Parser) ([]byte, *parser.LRTable ,error) {
    pa, err := parser.BuildSLRStates()
    if err != nil {
        return nil, nil,err
    }

    // Obtiene el JSON para el SLR completo
    slrJson, err := json.MarshalIndent(pa.SLR, "", "    ")
    if err != nil {
        return nil, nil,fmt.Errorf("error serializando SLR a JSON: %v", err)
    }

    // Obtiene el output de PrintSLR en formato JSON
    printSlrJson, err := pa.PrintSLR()
    if err != nil {
        return nil, nil,fmt.Errorf("error obteniendo la salida de PrintSLR: %v", err)
    }

    // Construye la tabla de análisis LR(1)
    lr1Table, err := parser.BuildLR1Table()
    if err != nil {
        return nil, nil,fmt.Errorf("error construyendo la tabla LR(1): %v", err)
    }

    // Serializa la tabla LR(1) a JSON
    lr1TableJson, err := json.MarshalIndent(lr1Table, "", "    ")

	// lr1Table.PrintTable()
    if err != nil {
        return nil, nil,fmt.Errorf("error serializando la tabla LR(1) a JSON: %v", err)
    }

    // Uniendo SLR, PrintSLR, y la tabla LR(1) en un solo objeto JSON
    type CombinedSLR struct {
        SLR       json.RawMessage `json:"slr"`
        PrintSLR  json.RawMessage `json:"print_slr"`
        LR1Table  json.RawMessage `json:"lr1_table"`
    }
    combined := CombinedSLR{
        SLR:      json.RawMessage(slrJson),
        PrintSLR: json.RawMessage(printSlrJson),
        LR1Table: json.RawMessage(lr1TableJson),
    }

    finalJson, err := json.MarshalIndent(combined, "", "    ")
    if err != nil {
        return nil, nil,fmt.Errorf("error combinando SLR, PrintSLR, y LR(1) en JSON: %v", err)
    }

    return finalJson,lr1Table, nil
}


func parseString(parser *parser.Parser, parsing string, parseTable *parser.LRTable) ([]byte, error) {
    validated, err := parser.ParseString(parsing, parseTable)

	response := utils.BoolsToBytes([]bool{validated})
	
    return response, err
}



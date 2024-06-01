package parser

import (
	// "fmt"
	// "sync"

	// "github.com/adrianfulla/compiler/backend/automatas"
	"fmt"

	"github.com/adrianfulla/compiler/backend/automatas"
	"github.com/adrianfulla/compiler/backend/utils"
)

type ParserScanner struct {
	Tokens      []utils.ParseToken
	IgnoredTokens []utils.ParseToken
	Productions []utils.ProductionToken
	SLR 		*automatas.SLR
}

func (parserScanner *ParserScanner) PrintParser(){
	fmt.Println("Accepted Tokens:")
	fmt.Println(parserScanner.Tokens)
	fmt.Println("Ignored Tokens:")
	fmt.Println(parserScanner.IgnoredTokens)
	fmt.Print("Productions: [\n")
	for _, prod := range parserScanner.Productions{
		fmt.Printf("\tHead: %s, Body %s\n",
		prod.Head, prod.Body)
	}
	fmt.Println("]")
}






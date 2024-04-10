package lexer

import (
	// "fmt"

	"fmt"

	"github.com/adrianfulla/compiler/backend/automatas"
	"github.com/adrianfulla/compiler/backend/utils"
)

type Lexer struct {
	file       string
	afdStack   utils.Stack
	TokenStack utils.Stack `json:"token_stack"`
}

func LexYmlFile(fileYml string) (*utils.Stack, error) {
	// fmt.Print(fileYml)
	lex := &Lexer{
		file: fileYml,
	}
	// automatas.ExtendedInfixToPosfix("'AB''\\t''\\''")
	automatas.ExtendedInfixToPosfix("'A''ABCD'\"abc\"")
	in, err := automatas.InfixToPosfix("abc(abc)(d)")
	if err == nil{
		fmt.Print(in)
	}
	// automatas.ExtendedInfixToPosfix("'AB''\\tABV'")
	// automatas.ExtendedInfixToPosfix("\"abc\\\\\\t\"")
	// automatas.ExtendedInfixToPosfix("A_B")
	// automatas.ExtendedInfixToPosfix("['A']")
	// automatas.ExtendedInfixToPosfix("['A'-'Z''a'-'z']")
	// automatas.ExtendedInfixToPosfix("[\"abcd\"]")
	// automatas.InfixToPosfix("[^'B'-'F']")
	// automatas.ExtendedInfixToPosfix("['A'-'C']#['B'-'F']")
	// automatas.ExtendedInfixToPosfix("(a)")
	// automatas.ExtendedInfixToPosfix("(['A'-'Z'])")

	stack, err := lex.parseFile()
	if err != nil {
		return nil, err
	}
	return &stack, nil
}

func (lex *Lexer) parseFile() (utils.Stack, error) {

	return lex.TokenStack, nil
}

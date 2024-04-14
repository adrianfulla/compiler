package lexer

import (
	// "fmt"

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
	// automatas.ExtendedInfixToPosfix("'ab\\'cd'")
	// automatas.ExtendedInfixToPosfix("\"\\\"\"")
	// automatas.ExtendedInfixToPosfix("'\\''\"\\\"\"")
	// in, err := automatas.InfixToPosfix("abc(abc)*(d)")
	// if err == nil{
	// 	fmt.Print(in)
	// }
	// automatas.ExtendedInfixToPosfix("'AB''\\tABV'")
	// automatas.ExtendedInfixToPosfix("'AB'")
	// automatas.ExtendedInfixToPosfix("\"abc\\\\\\t\"")
	// automatas.ExtendedInfixToPosfix("A_")
	// automatas.ExtendedInfixToPosfix("'ABCD''a'")
	// automatas.ExtendedInfixToPosfix("['A']['A'-'Z''a'-'z'][\"abc\"]")
	// automatas.ExtendedInfixToPosfix("['A']#['A'-'Z']")
	// automatas.ExtendedInfixToPosfix("('avc')\"avc\"*")
	// automatas.ExtendedInfixToPosfix("\"avc\"+*")
	// automatas.ExtendedInfixToPosfix("'a'\"a\"'abc'?'a'")
	// automatas.ExtendedInfixToPosfix("'a'|'b'*")
	automatas.ExtendedInfixToPosfix("(ident|nident)*")
	// automatas.ExtendedInfixToPosfix("[^\"abc\"]")
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

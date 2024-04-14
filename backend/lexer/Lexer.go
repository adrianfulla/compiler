package lexer

import (
	// "fmt"

	// "fmt"

	"github.com/adrianfulla/compiler/backend/automatas"
	"github.com/adrianfulla/compiler/backend/utils"
)

type Lexer struct {
	file       string
	afdStack   map[string]automatas.DAfdJson
	TokenStack utils.Stack `json:"token_stack"`
}

func LexYmlFile(fileYml string) (*utils.Stack, error) {
	// fmt.Print(fileYml)
	lex := &Lexer{
		file: fileYml,
	}

	definitions := map[string]string{}
	definitions["LETTER"] = "['A'-'Z''a'-'z']"
	definitions["NUMBER"] = "['0'-'9']"
	definitions["COMMENT"] = "'(* '_*' *)'"
	// definitions["COMMENT"] = "'(* '['A']*"
	validatedDefinitions := map[string]utils.DoublyLinkedList{}
	
	for token, def := range definitions{
		validated,err := automatas.ExtendedValidation(def)
		if err != nil{
			return nil, err
		}
		validatedDefinitions[token] = *validated
	}
	posfixDefinitions := map[string][]utils.RegexToken{}
	for token, def := range validatedDefinitions{
		posfix, err := automatas.ExtendedInfixToPosfix(def, validatedDefinitions)
		if err != nil{
			return nil, err
		}
		posfixDefinitions[token] = posfix
	}

	stack, err := lex.parseFile()
	if err != nil {
		return nil, err
	}
	return &stack, nil
}

func (lex *Lexer) parseFile() (utils.Stack, error) {

	return lex.TokenStack, nil
}

package lexer

import (
	// "fmt"

	// "fmt"

	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/adrianfulla/compiler/backend/automatas"
	"github.com/adrianfulla/compiler/backend/utils"
)

type Lexer struct {
	file       string
	afdStack   map[string]automatas.DAfdJson
	TokenStack utils.Stack `json:"token_stack"`
}

func LexYmlFile(fileYml string) (*Scanner, error) {
	// fmt.Print(fileYml)
	lex := &Lexer{
		file: fileYml,
		afdStack: make(map[string]automatas.DAfdJson),
	}

	definitions := map[string]string{}

	// definitions["COMMENT"] = "'(* '([\"ABCDEFGHIJKLMNOPQRSTUVWXYZ\"]*)' *)'"
	// definitions["SEMICOLON"] = "';'"
	definitions["COMMENTS"] = "'(* '['A'-'Z''a'-'z''0'-'9'\" .\"]*' *)'"
	definitions["DEFINITIONS"] = "'let '['A'-'Z''a'-'z']*\" = \"['A'-'Z''a'-'z''0'-'9'\"| []()\\'\\\"\\\\-*+?/%:;^.\"]*"
	definitions["TOKENRULES"] = "'rule tokens = '"
	// definitions["TOKENS"] = "'|'?'(\\'')|('\"')?[^\"\\t\\n\\r\\b\\f\\v\\\"]*"
	definitions["TOKENEXPRESIONS"] = "(['A'-'Z''a'-'z']+'\\n')|(\"'\"[^\"\\t\\s\\n \"]+\"'\")|(['A'-'Z''a'-'z']+[' ''\\t'][' ''\\t']+)"
	definitions["TOKENRETURNS"] = "\"{ return \"['A'-'Z']*\" }\""
	definitions["ERROR"]	= "_"
	
	validatedDefinitions := map[string]*utils.DoublyLinkedList{}
	
	for token, def := range definitions{
		validated,err := automatas.ExtendedValidation(def)
		if err != nil{
			fmt.Println(err)
			return nil, err
		}
		validatedDefinitions[token] = validated
	}
	// fmt.Println("Validated Def")
	posfixDefinitions := map[string][]utils.RegexToken{}
	for token, def := range validatedDefinitions{
		posfix, err := automatas.ExtendedInfixToPosfix(def, validatedDefinitions)
		if err != nil{
			fmt.Println(err)
			return nil, err
		}
		posfixDefinitions[token] = posfix
	}
	for token, posfix := range posfixDefinitions{
		afd := automatas.ExtendedNewDirectAfd(posfix)
		afdJson := afd.ToJson()
		// fmt.Println(afdJson)
		lex.afdStack[token] = *afdJson
	}

	Scanner, err := lex.parseFile()
	if err != nil {
		return nil, err
	}
	return &Scanner, nil
}

func (lex *Lexer) parseFile() (Scanner, error) {

	ch := make(chan map[string]utils.Stack, len(lex.afdStack))
	var wg sync.WaitGroup

	for index := range lex.afdStack{
		wg.Add(1)
		go func(el string){
			defer wg.Done()
			lex.searchYalex(el, ch)
		}(index)
	}

	wg.Wait()
	close(ch)

	tokensFound := []*automatas.AcceptedExp{}
	for maps := range ch {
		for index, result := range maps{
			// fmt.Println(index)
			for result.Size() > 0{
				res := result.Pop().(*automatas.AcceptedExp)
				res.Token = index
				// fmt.Println(res)
				tokensFound = AddOrUpdateExp(res, tokensFound)
			}
		}
	}
	tokensFound = SortTokens(tokensFound)
	tokensFoundStack := utils.Stack{}
	// fmt.Println("Tokens found")
	for _, tokenFound := range tokensFound{
		// fmt.Println(tokenFound)
		tokensFoundStack.Push(tokenFound)
	}

	newScanner := Scanner{
		Title: "",
		Header: []string{},
		Footer: []string{},
		Definitions: make(map[string]string),
		Tokens: []utils.LexToken{},
	}
	passedHeader := false
	for tokensFoundStack.Size() > 0{
		token := tokensFoundStack.Pop().(*automatas.AcceptedExp)
		// fmt.Printf("Case %s %s\n", token.Token, token.Value)
		switch token.Token{
		case "COMMENTS":
			// fmt.Printf("Case %s\n", token.Token)
			if newScanner.Title == ""{
				newScanner.Title = token.Value
			}else if !passedHeader {
				newScanner.Header = append(newScanner.Header, token.Value)
			}else{
				newScanner.Footer = append(newScanner.Header, token.Value)
			}
		case "DEFINITIONS":
			passedHeader = true
			// fmt.Printf("Case %s, %s\n", token.Token, token.Value)
			// println(token.Token)
			noLet := strings.TrimSpace(token.Value[4:])
			// fmt.Println(noLet)
			parts := strings.SplitN(noLet, "=", 2)
			if len(parts) != 2{
				fmt.Println("error parsing string")
				continue
			}
			if newScanner.Definitions[strings.TrimSpace(parts[0])] == ""{
				newScanner.Definitions[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
			}
			 
		case "TOKENEXPRESIONS":
			passedHeader = true
			// fmt.Printf("Case %s, %s\n", token.Token, token.Value)
			newLexToken := utils.LexToken{
				Regex: strings.TrimSpace(token.Value),
			}
			if tokensFoundStack.Peek().(*automatas.AcceptedExp).Token == "TOKENRETURNS"{
				token = tokensFoundStack.Pop().(*automatas.AcceptedExp)
				tokenName := token.Value[strings.Index(token.Value,"return")+6:strings.Index(token.Value,"}")]
				newLexToken.Token = strings.TrimSpace(tokenName)
				if !tokensFoundStack.IsEmpty() && tokensFoundStack.Peek().(*automatas.AcceptedExp).Token == "COMMENTS"{
					token = tokensFoundStack.Pop().(*automatas.AcceptedExp)
					if !tokensFoundStack.IsEmpty(){
						newLexToken.Action = strings.TrimSpace(token.Value)
					}else{
						tokensFoundStack.Push(token)
					}
				
				}
			}else{
				newLexToken.Token = strings.TrimSpace(token.Value)
			}
			if newScanner.Definitions[newLexToken.Token] == ""{

				newScanner.Definitions[newLexToken.Token] = newLexToken.Regex
			}
			newScanner.Tokens = append(newScanner.Tokens, newLexToken)
		case "TOKENRETURNS":
			passedHeader = true
			// fmt.Printf("Case %s\n", token.Token)
			// tokensFoundStack.Pop()
		default:
			// fmt.Printf("Default case %s\n", token.Token)
			// tokensFoundStack.Pop()
		}
	}
	return newScanner, nil
}

func SortTokens(tokens []*automatas.AcceptedExp) []*automatas.AcceptedExp{
	sort.Slice(tokens, func(i, j int) bool{
		return tokens[i].Start > tokens[j].Start
	})
	return tokens
}

func AddOrUpdateExp(newExp *automatas.AcceptedExp, currentExps []*automatas.AcceptedExp) []*automatas.AcceptedExp{
	temp := []*automatas.AcceptedExp{}
	for _, exp := range currentExps {
		if newExp.Start <= exp.Start && newExp.End >= exp.End{
			continue
		}
		if exp.Start <= newExp.Start && exp.End >= newExp.End{
			return currentExps
		}
		temp = append(temp, exp)
	}
	temp = append(temp, newExp)
	currentExps = temp
	return currentExps
}

func (lex *Lexer) searchYalex(index string, ch chan<-map[string]utils.Stack) {
	resultado := make(map[string]utils.Stack)
	resultado[index] = automatas.ExtendedSimulateAfd(lex.file, lex.afdStack[index])
	ch <- resultado
}



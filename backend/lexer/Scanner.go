package lexer

import (
	"fmt"
	"sync"

	"github.com/adrianfulla/compiler/backend/automatas"
	"github.com/adrianfulla/compiler/backend/utils"
)

type Scanner struct {
	Title       string            `json:"title"`
	Header      []string          `json:"header"`
	Footer      []string          `json:"footer"`
	Definitions map[string]string `json:"definitions"`
	Tokens      []utils.LexToken
}

func (scan *Scanner) ScanFile(file string) ([]*automatas.AcceptedExp, error) {
	// scan.PrintScanner()
	afdStack := map[string]automatas.DAfdJson{}
	validatedDefinitions := map[string]utils.DoublyLinkedList{}

	for token, def := range scan.Definitions {
		// fmt.Println(token)
		validated, err := automatas.ExtendedValidation(def)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		// fmt.Println(token)
		// validated.PrintForward()
		validatedDefinitions[token] = *validated
	}
	

	for _, token := range scan.Tokens {
		if validatedDefinitions[token.Token].Head != nil{
			continue
		}
		validated, err := automatas.ExtendedValidation(token.Regex)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		// fmt.Println(token)
		// validated.PrintForward()
		validatedDefinitions[token.Token] = *validated
	}

	// for token, def := range validatedDefinitions{
	// 	fmt.Print(token)
	// 	def.PrintForward()
	// }

	newDict, err := automatas.ReplaceReferenceIds(validatedDefinitions)
	if err != nil{
		return nil, err
	}
	validatedDefinitions = newDict
	
	tokenDefinitions := map[string]utils.DoublyLinkedList{}
	for _, token := range scan.Tokens{
		tokenDefinitions[token.Token] = validatedDefinitions[token.Token]
	}

	

	posfixDefinitions := map[string][]utils.RegexToken{}
	for token, def := range tokenDefinitions { 
		posfix, err := automatas.ExtendedInfixToPosfix(def, validatedDefinitions)	
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		// fmt.Println(token, posfix)
		posfixDefinitions[token] = posfix
	}
	
	for token, posfix := range posfixDefinitions {
		afd := automatas.ExtendedNewDirectAfd(posfix)
		afdJson := afd.ToJson()
		// fmt.Println(afdJson)
		afdStack[token] = *afdJson
	}

	return scan.parseFile(file, afdStack)
}

func (scan *Scanner) parseFile(file string, afdStack map[string]automatas.DAfdJson) ([]*automatas.AcceptedExp, error) {
	// fmt.Println(file)
	ch := make(chan map[string]utils.Stack, len(afdStack))
	var wg sync.WaitGroup
	for index := range afdStack {
		wg.Add(1)
		go func(el string) {
			defer wg.Done()
			scan.searhFile(el, afdStack[el], file, ch)
		}(index)
	}
	wg.Wait()
	close(ch)
	tokensFound := []*automatas.AcceptedExp{}
	for maps := range ch {
		for index, result := range maps {
			fmt.Println(index)
			for result.Size() > 0 {
				res := result.Pop().(*automatas.AcceptedExp)
				res.Token = index
				fmt.Println(res)
				tokensFound = AddOrUpdateExp(res, tokensFound)
			}
		}
	}
	tokensFound = SortTokens(tokensFound)
	return tokensFound, nil
}

func (scan *Scanner) searhFile(index string, afd automatas.DAfdJson, file string, ch chan<- map[string]utils.Stack) {
	resultado := make(map[string]utils.Stack)
	resultado[index] = automatas.ExtendedSimulateAfd(file, afd)
	ch <- resultado
}

func (scan *Scanner) PrintScanner(){
	fmt.Printf("Title for Scanner: %s\n", scan.Title)
	fmt.Printf("Header for Scanner: %s\n", scan.Header)
	fmt.Printf("Definitions for Scanner:\n" )
	for token, def := range scan.Definitions{
		fmt.Printf("Defition of %s is %s\n",token, def)
	}
	fmt.Printf("Tokens for Scanner:\n" )
	for _,token := range scan.Tokens{
		fmt.Printf("Token %s with regex %s and action %s\n",token.Token,token.Regex, token.Action )
	}

	fmt.Printf("Footer for Scanner: %s\n", scan.Footer)
}
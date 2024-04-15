package lexer

import (
	"github.com/adrianfulla/compiler/backend/automatas"
	"github.com/adrianfulla/compiler/backend/utils"
	"fmt"
	"sync"
)

type Scanner struct {
	Title string `json:"title"`
	Header []string	`json:"header"`
	Footer []string `json:"footer"`
	Definitions map[string]string `json:"definitions"`
	Tokens []utils.LexToken
}


func (scan *Scanner) ScanFile(file string) ([]*automatas.AcceptedExp, error){
	afdStack := map[string]automatas.DAfdJson{}
	validatedDefinitions := map[string]utils.DoublyLinkedList{}
	
	for token, def := range scan.Definitions{
		fmt.Println(token)
		validated,err := automatas.ExtendedValidation(def)
		if err != nil{
			fmt.Println(err)
			return nil, err
		}
		validatedDefinitions[token] = *validated
	}
	for _, token := range scan.Tokens{
		fmt.Println(token)
		validated,err := automatas.ExtendedValidation(token.Regex)
		if err != nil{
			fmt.Println(err)
			return nil, err
		}
		validatedDefinitions[token.Token] = *validated
	}
	posfixDefinitions := map[string][]utils.RegexToken{}
	for token, def := range validatedDefinitions{
		posfix, err := automatas.ExtendedInfixToPosfix(def, validatedDefinitions)
		if err != nil{
			fmt.Println(err)
			return nil, err
		}
		// fmt.Println(token,posfix)
		posfixDefinitions[token] = posfix
	}
	for token, posfix := range posfixDefinitions{
		afd := automatas.ExtendedNewDirectAfd(posfix)
		afdJson := afd.ToJson()
		// fmt.Println(afdJson)
		afdStack[token] = *afdJson
	}

	return scan.parseFile(file, afdStack)
}

func (scan *Scanner) parseFile(file string, afdStack map[string]automatas.DAfdJson) ([]*automatas.AcceptedExp, error){
	ch := make(chan map[string]utils.Stack, len(afdStack))
	var wg sync.WaitGroup
	for index := range afdStack{
		wg.Add(1)
		go func(el string){
			defer wg.Done()
			scan.searhFile(el, afdStack[el],file, ch)
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
	return tokensFound, nil
}

func (scan *Scanner) searhFile(index string,afd automatas.DAfdJson, file string, ch chan<-map[string]utils.Stack){
	resultado := make(map[string]utils.Stack)
	resultado[index] = automatas.ExtendedSimulateAfd(file, afd)
	ch <- resultado
}
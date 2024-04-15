package automatas

import (
	"fmt"
	"strconv"
	"strings"
	// "time"

	"unicode"

	"github.com/adrianfulla/compiler/backend/utils"
)

func validation(regex string) (string, error) {
	if len(regex) == 0 {
		return "", fmt.Errorf("ingrese una expresión regular")
	}

	if strings.Contains(regex, "^") || strings.Contains(regex, ".") {
		return "", fmt.Errorf("el símbolo de concatenación se agregará automáticamente después. Ingrese su expresión sin el símbolo de concatenación")
	}

	for _, c := range regex {
		if !unicode.IsLetter(c) && !unicode.IsNumber(c) && !strings.ContainsRune("*+?|()ε", c) {
			return "", fmt.Errorf("carácter no válido en la expresión regular")
		}
	}

	stack := []rune{}
	for _, char := range regex {
		switch char {
		case '(':
			stack = append(stack, char)
		case ')':
			if len(stack) == 0 {
				return "", fmt.Errorf("paréntesis no balanceados en la expresión regular")
			}
			stack = stack[:len(stack)-1] // Simula el pop
		}
	}

	if len(stack) > 0 {
		return "", fmt.Errorf("paréntesis no balanceados en la expresión regular")
	}

	operators := "*+?|"
	for i := 0; i < len(regex)-1; i++ {
		currentChar := regex[i]
		nextChar := regex[i+1]

		if strings.ContainsRune(operators, rune(currentChar)) && strings.ContainsRune(operators, rune(nextChar)) && nextChar != '|' {
			return "", fmt.Errorf("sintaxis incorrecta de operadores en la expresión regular")
		}
	}
	return regex, nil
}

func escapeRune(char rune) (rune, error) {
	// Crear la secuencia de escape como una cadena
	escapeSequence := fmt.Sprintf("\\%c", char)
	// Interpretar la secuencia de escape
	decodedString, err := strconv.Unquote(`"` + escapeSequence + `"`)
	if err != nil {
		return 0, err
	}
	// Como sabemos que es una secuencia de escape, debe tener un solo carácter
	return []rune(decodedString)[0], nil
}

func addBetweenOperator(nodo *utils.LinkedNode, char utils.RegexToken) *utils.LinkedNode {
	newPrevToken := utils.LinkedNode{
		Value: char,
		Next:  nodo,
		Prev:  nodo.Prev,
	}
	nodo.Prev.Next = &newPrevToken
	nodo.Prev = &newPrevToken

	return nodo
}

func shouldConcat(nodo *utils.LinkedNode) bool {
	prev := nodo.Prev.Value.(utils.RegexToken)
	return (prev.IsOperator != "OROPERATOR" &&
		prev.IsOperator != "OPENPARENTHESES" &&
		nodo.Value.(utils.RegexToken).IsOperator != "CLOSEPARENTHESES" &&
		nodo.Value.(utils.RegexToken).IsOperator != "KLEENEOPERATOR")
}

func addOpenParentheses(nodo *utils.LinkedNode, list *utils.DoublyLinkedList) *utils.LinkedNode {
	prevCharToken := utils.RegexToken{
		Value:      []rune{'('},
		IsOperator: "OPENPARENTHESES",
	}
	if nodo.Prev == nil {
		list.Head = nodo.Next
		list.Prepend(prevCharToken)
		nodo = list.Head
	} else{
		// fmt.Print("Prev nodo found")
		// fmt.Print(nodo.Prev.Value)
		// fmt.Print("\n")
		isConcat :=shouldConcat(nodo)
		if  isConcat{
			catChar := utils.RegexToken{
				Value:      []rune{'^'},
				IsOperator: "CATOPERATOR",
			}
			nodo = addBetweenOperator(nodo, catChar)
		}
		// list.PrintForward()
			// fmt.Print("IN Add PAR\n")
			// fmt.Printf("Prev Token %s \n", nodo.Prev.Value)
			// fmt.Printf("Current Token %s \n", nodo.Value)
			// fmt.Printf("Next Token %s \n", nodo.Next.Value)
			// fmt.Printf("Next Token %s \n\n", nodo.Next.Next.Value)
		newPrevToken := utils.LinkedNode{
			Value: prevCharToken,
			Next:  nodo.Next,
			Prev:  nodo.Prev,
		}
		nodo.Next.Prev = &newPrevToken
		nodo.Prev.Next = &newPrevToken
		if isConcat{
			return nodo.Prev.Prev.Next.Next
		}
			// fmt.Printf("Prev Token %s \n", nodo.Prev.Value)
			// fmt.Printf("Prev Next Token %s \n", nodo.Prev.Next.Value)
			// fmt.Printf("Current Token %s \n", nodo.Value)
			// fmt.Printf("Next Token %s \n", nodo.Next.Value)
			// fmt.Printf("Next Token %s \n\n", nodo.Next.Next.Value)
			// fmt.Printf("Next Prev Token %s \n\n", nodo.Next.Prev.Value)
		nodo = nodo.Prev.Next
	}
	// fmt.Print("\n\n")
	// list.PrintForward()
	return nodo
}

func addCloseParentheses(nodo *utils.LinkedNode, list *utils.DoublyLinkedList) *utils.LinkedNode {
	nextChar := utils.RegexToken{
		Value:      []rune{')'},
		IsOperator: "CLOSEPARENTHESES",
	}
	// fmt.Print("ADD CLOSEPAR")
	// fmt.Print(nodo.Value)
	// fmt.Print("\n")
	if nodo.Next == nil {
		// fmt.Println("Encontro null")
		// fmt.Printf("Current Token %s \n", currentToken.Value)
		// fmt.Printf("Next Token %s \n", currentToken.Next.Value)
		list.Tail = nodo.Prev
		list.Append(nextChar)
	} else {
		newNextToken := utils.LinkedNode{
			Value: nextChar,
			Next:  nodo.Next,
			Prev:  nodo.Prev,
		}
		nodo.Prev.Next = &newNextToken
		nodo.Next.Prev = &newNextToken
	}

	return nodo
}

func findEarliestAcceptedParentheses(nodo *utils.LinkedNode) (*utils.LinkedNode){
	deepness := 0
	for nodo.Prev != nil && (nodo.Value.(utils.RegexToken).IsOperator != "OPENPARENTHESES" || deepness != 0){
		if nodo.Value.(utils.RegexToken).IsOperator == "CLOSEPARENTHESES"{
			deepness ++
		}else if nodo.Value.(utils.RegexToken).IsOperator == "OPENPARENTHESES"{
			deepness --
		}
		if deepness == 0{
			break
		}
		nodo = nodo.Prev
	}


	return nodo
}

func ExtendedValidation(regex string) (*utils.DoublyLinkedList,error) {
	operators := []string{
		"OPENPARENTHESES",
		"CLOSEPARENTHESES",
		"KLEENE",
		"CATOPERATOR",
		"OROPERATOR",
	}
	fmt.Println(regex)
	if len(regex) == 0 {
		return nil, fmt.Errorf("regex error: empty regular expression")
	}
	regexDLinkedList := utils.DoublyLinkedList{}

	for _, char := range regex {
		regexDLinkedList.Append(utils.RegexToken{
			Value: []rune{char},
		})
	}
	currentToken := regexDLinkedList.Head
	openPar := 0
	isDiff := false
	firstCharSet := []rune{}
	for currentToken != nil {
		fmt.Println("")
		fmt.Println("Forward")
		regexDLinkedList.PrintForward()
		fmt.Println("Reverse")
		regexDLinkedList.PrintReverse()
		char := currentToken.Value.(utils.RegexToken).Value
		if len(char) == 1 {
			char := char[len(char)-1]
			fmt.Printf("Encontrado:%s\n", string(char))
			switch char {
			case '\'':
				fmt.Printf("Case:%s\n", string(char))
				if currentToken.Next == nil {
					
					return nil, fmt.Errorf("regex error: regex parsing error, found extra '")
				}
				
				currentToken = addOpenParentheses(currentToken, &regexDLinkedList)
				
				currentToken = currentToken.Next
				
				catChar := utils.RegexToken{
					Value:      []rune{'^'},
					IsOperator: "CATOPERATOR",
				}
				
				for currentToken.Next.Value.(utils.RegexToken).Value[0] != '\'' || currentToken.Value.(utils.RegexToken).Value[0] == '\\' {
					fmt.Println("ACA")
					if currentToken.Value.(utils.RegexToken).Value[0] == '\\' {
						if currentToken.Next.Value.(utils.RegexToken).Value[0] != '\'' {
							escaped, err := escapeRune(currentToken.Next.Value.(utils.RegexToken).Value[0])
							if err != nil {
								escaped = currentToken.Next.Next.Value.(utils.RegexToken).Value[0]
							}
							currentToken.Value = utils.RegexToken{
								Value: []rune{escaped},
							}
							currentToken.Next = currentToken.Next.Next
							currentToken.Next.Prev = currentToken
						} else {
							fmt.Println("ACA")
							currentToken.Prev.Next = currentToken.Next
							currentToken.Next.Prev = currentToken.Prev
							currentToken = currentToken.Next
						}

					}

					if currentToken.Next != nil && currentToken.Next.Value.(utils.RegexToken).Value[0] != '\'' {
						currentToken = addBetweenOperator(currentToken.Next, catChar)
					}
				}
				// fmt.Println("ERROR ACA\n")
				currentToken = addCloseParentheses(currentToken.Next, &regexDLinkedList)
				
			case '"':
				fmt.Printf("Case:%s\n", string(char))
				if currentToken.Next == nil {
					return nil, fmt.Errorf("regex error: regex parsing error, found extra '")
				}

				currentToken = addOpenParentheses(currentToken, &regexDLinkedList)

				currentToken = currentToken.Next

				catChar := utils.RegexToken{
					Value:      []rune{'^'},
					IsOperator: "CATOPERATOR",
				}
				for currentToken.Next.Value.(utils.RegexToken).Value[0] != '"' || currentToken.Value.(utils.RegexToken).Value[0] == '\\' {
					if currentToken.Value.(utils.RegexToken).Value[0] == '\\' {
						if currentToken.Next.Value.(utils.RegexToken).Value[0] != '"' {
							escaped, err := escapeRune(currentToken.Next.Value.(utils.RegexToken).Value[0])
							if err != nil {
								escaped = currentToken.Next.Next.Value.(utils.RegexToken).Value[0]
							}
							currentToken.Value = utils.RegexToken{
								Value: []rune{escaped},
							}
							currentToken.Next = currentToken.Next.Next
							currentToken.Next.Prev = currentToken
						} else {
							currentToken.Prev.Next = currentToken.Next
							currentToken.Next.Prev = currentToken.Prev
							currentToken = currentToken.Next
						}

					}

					if currentToken.Next != nil && currentToken.Next.Value.(utils.RegexToken).Value[0] != '"' {
						currentToken = addBetweenOperator(currentToken.Next, catChar)
					}
				}

				currentToken = addCloseParentheses(currentToken.Next, &regexDLinkedList)
			case '_':
				fmt.Printf("Case:%s\n", string(char))
				
				nextToken := currentToken.Next
				currentToken = &utils.LinkedNode{
					Value: utils.RegexToken{},
					Prev:  currentToken.Prev,
					Next:  currentToken,
				}

				currentToken = addOpenParentheses(currentToken, &regexDLinkedList).Next
				orChar := utils.RegexToken{
					Value:      []rune{'|'},
					IsOperator: "OROPERATOR",
				}
				isLast := currentToken.Next == nil

				for i := 0; i <= 254; i++ {
					currentToken.Value = utils.RegexToken{
						Value: []rune{rune(i)},
					}

					currentToken.Next = &utils.LinkedNode{
						Value: orChar,
						Next: &utils.LinkedNode{
							Prev: currentToken.Next,
						},
						Prev: currentToken,
					}
					currentToken.Next.Next.Prev = currentToken.Next
					currentToken = currentToken.Next.Next
					if isLast {
						regexDLinkedList.Tail = currentToken
					}

				}
				// fmt.Printf("Prev Token %s \n", currentToken.Prev.Value)
				// fmt.Printf("Current Token %s \n", currentToken.Value)
				// fmt.Printf("Next Token %s \n", currentToken.Next.Value)

				currentToken.Value = utils.RegexToken{
					Value: []rune{rune(255)},
				}
				// fmt.Printf("\nPrev Token %s \n", currentToken.Prev.Value)
				// fmt.Printf("Current Token %s \n", currentToken.Value)
				// fmt.Printf("Next Token %s \n", currentToken.Next.Value)
				emptyToken := &utils.LinkedNode{
					Value: utils.RegexToken{},
					Prev: currentToken,
					Next: nextToken,
				}

				currentToken = addCloseParentheses(emptyToken, &regexDLinkedList)
				// fmt.Printf("\nPrev Token %s \n", currentToken.Prev.Value)
				// fmt.Printf("Current Token %s \n", currentToken.Value)
				// fmt.Printf("Next Token %s \n", currentToken.Next.Value)
			case '[':
				fmt.Printf("Case:%s\n", string(char))

				currentToken = addOpenParentheses(currentToken, &regexDLinkedList)
				tempToken := currentToken.Next
				notInCharset := false
				symbolsStack := utils.Stack{}
				for tempToken.Value.(utils.RegexToken).Value[0] != ']' {
					switch tempToken.Value.(utils.RegexToken).Value[0] {
					case '\'':
						// fmt.Print("Case ' found\n")
						// fmt.Printf("\nPrev Token %s \n", tempToken.Prev.Value)
						// fmt.Printf("Current Token %s \n", tempToken.Value)
						// fmt.Printf("Next Token %s \n", tempToken.Next.Value)
						if tempToken.Next.Value.(utils.RegexToken).Value[0] =='\\'{
							escaped, err :=escapeRune(tempToken.Next.Next.Value.(utils.RegexToken).Value[0])
							if err != nil{
								escaped = tempToken.Next.Next.Value.(utils.RegexToken).Value[0]
							}
							tempToken.Next.Next.Value = utils.RegexToken{
								Value: []rune{escaped},
							}
							tempToken = tempToken.Next
						}
						symbolsStack.Push(tempToken.Next.Value.(utils.RegexToken).Value[0])
						if tempToken.Next.Next.Next.Value.(utils.RegexToken).Value[0] == '-' {
							tempToken = tempToken.Next.Next.Next.Next
							firstLimit := symbolsStack.Peek().(rune)
							secondLimit := tempToken.Next.Value.(utils.RegexToken).Value[0]
							// fmt.Printf("First: %s, Second: %s\n", firstLimit, secondLimit)
							tempstack, err := expandBrackets(firstLimit+1, secondLimit)
							if err != nil {
								return nil, fmt.Errorf("invalid character set")
							}
							for tempstack.Size() > 0 {
								symbolsStack.Push(tempstack.Pop())
							}

						}
						tempToken = tempToken.Next.Next.Next
					case '"':
						if tempToken.Next.Value.(utils.RegexToken).Value[0] =='\\'{
							escaped, err :=escapeRune(tempToken.Next.Next.Value.(utils.RegexToken).Value[0])
							if err != nil{
								escaped = tempToken.Next.Next.Value.(utils.RegexToken).Value[0]
							}
							tempToken.Next.Next.Value = utils.RegexToken{
								Value: []rune{escaped},
							}
							tempToken = tempToken.Next
						}
						tempToken = tempToken.Next
						for tempToken.Value.(utils.RegexToken).Value[0] != '"' {
							symbolsStack.Push(tempToken.Value.(utils.RegexToken).Value[0])
							tempToken = tempToken.Next
						}
						tempToken = tempToken.Next
					case '^':
						// fmt.Print("Case ^ found\n")
						notInCharset = true
						tempToken = tempToken.Next
					default:
						// fmt.Print("Case nil found\n")
						return nil, fmt.Errorf("invalid character set")
					}

				}
				if notInCharset {
					notIn := []rune{}
					for symbolsStack.Size() > 0 {
						notIn = append(notIn, symbolsStack.Pop().(rune))
					}
					// fmt.Printf("NOT IN %s\n", string(notIn))
					for i := range 255 {
						if !strings.ContainsRune(string(notIn), rune(i)) {
							symbolsStack.Push(rune(i))
						}
					}
				}
				if !isDiff{
					firstElem := symbolsStack.Peek().(rune)
					firstCharSet = []rune{}
					firstCharSet = append(firstCharSet, firstElem)
				}else{
					secondCharSet := []rune{}
					tempStackSymbol := utils.Stack{}
					for symbolsStack.Size() > 0 {
						char :=symbolsStack.Pop().(rune)
						secondCharSet = append(secondCharSet, char)
						if !strings.ContainsRune(string(firstCharSet), char){
							tempStackSymbol.Push(char)
						}
					}
					for _, char := range firstCharSet{
						if !strings.ContainsRune(string(secondCharSet), char){
							tempStackSymbol.Push(char)
						}
					}
					for tempStackSymbol.Size() > 0{
						symbolsStack.Push(tempStackSymbol.Pop())
					}
					for currentToken.Prev.Value.(utils.RegexToken).IsOperator != "OPENPARENTHESES"{
						currentToken = currentToken.Prev
					}
					currentToken = currentToken.Prev
				}
				newNodo := utils.LinkedNode{
					Value: utils.RegexToken{
						Value: []rune{symbolsStack.Pop().(rune)},
					},
					Prev: currentToken,
					Next: tempToken,
				}
				orChar := utils.RegexToken{
					Value:      []rune{'|'},
					IsOperator: "OROPERATOR",
				}
				currentToken.Next = &newNodo
				tempToken.Prev = &newNodo
				currentToken = currentToken.Next
				for symbolsStack.Size() > 0 {
					charElem :=symbolsStack.Pop().(rune)
					// fmt.Printf("Char in symbolstack: %s\n", charElem)
					firstCharSet = append(firstCharSet, charElem)
					newNodo := utils.LinkedNode{
						Value: utils.RegexToken{
							Value: []rune{charElem},
						},
						Prev: currentToken,
						Next: tempToken,
					}
					currentToken.Next = &newNodo
					tempToken.Prev = &newNodo
					currentToken = addBetweenOperator(&newNodo, orChar)
				}
				currentToken = currentToken.Next
				
				// fmt.Printf("Prev Token %s \n", currentToken.Prev.Value)
				// fmt.Printf("Current Token %s \n", currentToken.Value)
				// fmt.Printf("Next Token %s \n", currentToken.Next.Value)
				currentToken = addCloseParentheses(currentToken, &regexDLinkedList)
			case ']':
				return nil, fmt.Errorf("unexpected closing bracket found at regex")
			case '(':
				fmt.Printf("Case:%s\n", string(char))
				currentToken.Value = utils.RegexToken{
					Value: currentToken.Value.(utils.RegexToken).Value,
					IsOperator: "OPENPARENTHESES",
				}
				openPar ++
			case ')':
				fmt.Printf("Case:%s\n", string(char))
				currentToken.Value = utils.RegexToken{
					Value: currentToken.Value.(utils.RegexToken).Value,
					IsOperator: "CLOSEPARENTHESES",
				}
				openPar --
			case '^':
				fmt.Printf("Case:%s\n", string(char))
				return nil, fmt.Errorf("regex error: character %s not valid in regex", string(char))
			case '#':
				fmt.Printf("Case:%s\n", string(char))
				isDiff = true
			case '*':
				fmt.Printf("Case:%s\n", string(char))
				// fmt.Printf("Prev Token %s \n", currentToken.Prev.Value)
				// fmt.Printf("Current Token %s \n", currentToken.Value)
				// fmt.Printf("Next Token %s \n", currentToken.Next.Value)
				tempToken := currentToken.Prev
				tempToken = findEarliestAcceptedParentheses(tempToken)
				// fmt.Print("After earliest\n")
				// fmt.Printf("Prev Token %s \n", tempToken.Prev.Value)
				// fmt.Printf("Current Token %s \n", tempToken.Value)
				// fmt.Printf("Next Token %s \n", tempToken.Next.Value)
				if tempToken.Prev != nil{
					tempToken = &utils.LinkedNode{
						Value: utils.RegexToken{
							Value: []rune{'('},
							IsOperator: "OPENPARENTHESES",
						} ,
						Prev: tempToken.Prev,
						Next: tempToken,
					}
					tempToken.Prev.Next = tempToken
					tempToken.Next.Prev = tempToken
				}else{
					regexDLinkedList.Prepend(utils.RegexToken{
						Value: []rune{'('},
						IsOperator: "OPENPARENTHESES",
					})
				}
				
				// addOpenParentheses(tempToken, &regexDLinkedList)
				
				emptyToken := &utils.LinkedNode{
					Value: utils.RegexToken{},
					Prev: currentToken,
					Next: currentToken.Next,
				}
				currentToken = addCloseParentheses(emptyToken, &regexDLinkedList).Prev
				currentToken.Value =	utils.RegexToken{
							Value: currentToken.Value.(utils.RegexToken).Value,
							IsOperator: "KLEENE",
						}
				currentToken = currentToken.Next
			case '+':
				fmt.Printf("Case:%s\n", string(char))
				tempToken := currentToken.Prev
				tempToken = findEarliestAcceptedParentheses(tempToken)


				firstEmptyToken := utils.LinkedNode{
					Value: utils.RegexToken{},
					Prev: tempToken,
					Next: tempToken.Next,
				}
				
				
				tempToken = addOpenParentheses(&firstEmptyToken,&regexDLinkedList)
				// regexDLinkedList.PrintForward()

				emptyToken := utils.LinkedNode{
					Value: utils.RegexToken{},
					Prev: currentToken.Prev,
					Next: currentToken,
				}
				currentToken = addOpenParentheses(&emptyToken,&regexDLinkedList).Next
				
				deepness := 0
				for tempToken.Next != nil && (tempToken.Value.(utils.RegexToken).IsOperator != "CLOSEPARENTHESES" || deepness != 0){
					
					if tempToken.Value.(utils.RegexToken).IsOperator == "CLOSEPARENTHESES"{
						deepness --
						
					}else if tempToken.Value.(utils.RegexToken).IsOperator == "OPENPARENTHESES"{
						deepness ++
					}
					currentToken.Prev = &utils.LinkedNode{
						Value: tempToken.Value,
						Prev: currentToken.Prev,
						Next: currentToken,
					}
					currentToken.Prev.Prev.Next = currentToken.Prev

					// fmt.Print("\nlooking for RPar\nforward\n")
					// fmt.Printf("toke %s, deep %d\n",tempToken.Value, deepness)
					// fmt.Print("\n")
					// regexDLinkedList.PrintForward()
					// fmt.Print("\n\nreverse\n")
					// regexDLinkedList.PrintReverse()
					// time.Sleep(1 * time.Second)
					tempToken = tempToken.Next

					if deepness == 0{
						break
					}
				}				


				// fmt.Printf("Next Token %s \n", currentToken.Next.Value)
				// fmt.Print("IN *\nforward\n")
				// regexDLinkedList.PrintForward()
				// fmt.Print("\n\nreverse\n")
				// regexDLinkedList.PrintReverse()
				emptyToken = utils.LinkedNode{
					Value: utils.RegexToken{},
					Prev: currentToken,
					Next: currentToken.Next,
				}

				currentToken = addCloseParentheses(&emptyToken, &regexDLinkedList).Prev
				currentToken.Value =	utils.RegexToken{
							Value: []rune{'*'},
							IsOperator: "KLEENE",
						}
				emptyToken = utils.LinkedNode{
					Value: utils.RegexToken{},
					Prev: currentToken,
					Next: currentToken.Next,
				}
				currentToken = addCloseParentheses(&emptyToken, &regexDLinkedList)
				currentToken = currentToken.Next
			case '?':
				fmt.Printf("Case:%s\n", string(char))
				// fmt.Print("Bef earlies\n")
				// fmt.Printf("Prev Token %s \n", currentToken.Prev.Value)
				// fmt.Printf("Current Token %s \n", currentToken.Value)
				// fmt.Printf("Next Token %s \n", currentToken.Next.Value)
				tempToken := currentToken.Prev
				tempToken = findEarliestAcceptedParentheses(tempToken)

				tempToken = &utils.LinkedNode{
					Value: utils.RegexToken{},
					Prev: tempToken,
					Next: tempToken.Next,
				}
				
				// fmt.Printf("Prev Temp Token %s \n", tempToken.Prev.Value)
				addOpenParentheses(tempToken, &regexDLinkedList)
				// regexDLinkedList.PrintForward()
				
				emptyToken := &utils.LinkedNode{
					Value: utils.RegexToken{},
					Prev: currentToken,
					Next: currentToken.Next,
				}
				currentToken = addCloseParentheses(emptyToken, &regexDLinkedList).Prev
				currentToken.Value =	utils.RegexToken{
							Value: []rune{},
							IsOperator: "NULL",
						}
				orChar := utils.RegexToken{
					Value:      []rune{'|'},
					IsOperator: "OROPERATOR",
				}
				currentToken = addBetweenOperator(currentToken, orChar)
				currentToken = currentToken.Next
			case '|':
				// fmt.Printf("Case:%s\n", string(char))
				currentToken = &utils.LinkedNode{
					Value: utils.RegexToken{
						Value: currentToken.Value.(utils.RegexToken).Value,
						IsOperator: "OROPERATOR",
					},
					Next: currentToken.Next,
					Prev: currentToken.Prev,
				}
				currentToken.Next.Prev = currentToken
				currentToken.Prev.Next = currentToken
			default:
				fmt.Printf("Default Case:%s\n", string(char))
				if currentToken.Prev == nil ||  utils.StringInStringArray(currentToken.Prev.Value.(utils.RegexToken).IsOperator, operators) {
					tempToken := &utils.LinkedNode{
						Value: utils.RegexToken{},
						Prev: currentToken.Prev,
						Next: currentToken,
					}
					// fmt.Printf("Prev Temp Token %s \n", tempToken.Prev.Value)
					currentToken = addOpenParentheses(tempToken, &regexDLinkedList).Next
					// regexDLinkedList.PrintForward()
					currentToken.Value = utils.RegexToken{
							Value: currentToken.Value.(utils.RegexToken).Value,
							IsOperator: string(currentToken.Value.(utils.RegexToken).Value),
					}
					currentToken.Prev.Next = currentToken
					// fmt.Printf("Prev Token %s \n", currentToken.Prev.Value)
					// fmt.Printf("Current Token %s \n", currentToken.Value)
					// fmt.Printf("Next Token %s \n", currentToken.Next.Value)
				}else if !utils.StringInStringArray(currentToken.Prev.Value.(utils.RegexToken).IsOperator, operators) {
					// fmt.Printf("Prev Token %s \n", currentToken.Prev.Value)
					// fmt.Printf("Current Token %s \n", currentToken.Value)
					// fmt.Printf("Next Token %s \n", currentToken.Next.Value)
					currentToken = &utils.LinkedNode{
						Value: utils.RegexToken{
							Value: append(currentToken.Prev.Value.(utils.RegexToken).Value, currentToken.Value.(utils.RegexToken).Value...),
							IsOperator: currentToken.Prev.Value.(utils.RegexToken).IsOperator + string(currentToken.Value.(utils.RegexToken).Value),
						},
						Next: currentToken.Next,
						Prev: currentToken.Prev.Prev,
					}
					currentToken.Prev.Next = currentToken
					if currentToken.Next != nil{
						currentToken.Next.Prev = currentToken
						currentToken = currentToken.Next.Prev
					}else{
						regexDLinkedList.Tail = currentToken
					}
					
				}else{
					currentToken.Value = utils.RegexToken{
						Value: currentToken.Value.(utils.RegexToken).Value,
						IsOperator: string(currentToken.Value.(utils.RegexToken).Value),
					}
				}
				if currentToken.Next == nil || strings.ContainsRune("'\"[]()#*|+", currentToken.Next.Value.(utils.RegexToken).Value[0]){
					// fmt.Printf("Prev Token %s \n", currentToken.Prev.Value)
					// fmt.Printf("Current Token %s \n", currentToken.Value)
					// fmt.Printf("Next Token %s \n", currentToken.Next.Value)
					emptyToken := &utils.LinkedNode{
						Value: utils.RegexToken{},
						Prev: currentToken,
						Next: currentToken.Next,
					}
					currentToken = addCloseParentheses(emptyToken, &regexDLinkedList)
				}
			}
			currentToken = currentToken.Next
		}
	}

	// fmt.Print("END OF SHIFT\n")
	// fmt.Print("forward\n")
	// regexDLinkedList.PrintForward()
	// fmt.Print("\n\nreverse\n")
	// regexDLinkedList.PrintReverse()
	if openPar != 0{
		return nil, fmt.Errorf("regex error: unbalanced parentheses")
	}
	return  &regexDLinkedList, nil
}

func expandBrackets(start rune, end rune) (utils.Stack, error) {
	result := utils.Stack{}
	for c := start; c <= end; c++ {
		result.Push(c)
	}

	return result, nil
}

func convertQuestionMarkAndPlusSign(regex string) (string, error) {
	var result strings.Builder
	i := 0

	for i < len(regex) {
		char := regex[i]

		if char == '?' {
			if i == 0 {
				return "", fmt.Errorf("uso incorrecto del '?' al principio de la expresión")
			}

			if regex[i-1] == ')' {
				// Encuentra el paréntesis abierto correspondiente
				balance := 1
				j := i - 2
				for j >= 0 && balance > 0 {
					if regex[j] == ')' {
						balance++
					} else if regex[j] == '(' {
						balance--
					}
					j--
				}

				if balance != 0 {
					return "", fmt.Errorf("paréntesis no balanceados")
				}

				// Añade el grupo con |ε
				result.Reset()
				result.WriteString(regex[:j+1] + "(" + regex[j+1:i] + "|ε)")
			} else {
				// Caso simple: ? se aplica directamente al caracter anterior
				prevChar := regex[i-1]
				result.WriteString("(" + string(prevChar) + "|ε)")
				regex = regex[i+1:]
				i = -1 // Resetear índice para continuar con el resto de la cadena
			}
		} else if char == '+' {
			if i == 0 {
				return "", fmt.Errorf("uso incorrecto del '+' al principio de la expresión")
			}

			if regex[i-1] == ')' {
				// Encuentra el paréntesis abierto correspondiente
				balance := 1
				j := i - 2
				for j >= 0 && balance > 0 {
					if regex[j] == ')' {
						balance++
					} else if regex[j] == '(' {
						balance--
					}
					j--
				}

				if balance != 0 {
					return "", fmt.Errorf("paréntesis no balanceados")
				}

				// Añade el grupo con *
				result.Reset()
				result.WriteString(regex[:j+1] + regex[j+1:i] + regex[j+1:i] + "*")
			} else {
				// Caso simple: ? se aplica directamente al caracter anterior
				prevChar := regex[i-1]
				result.WriteString(string(prevChar) + "*")
				regex = regex[i+1:]
				i = -1 // Resetear índice para continuar con el resto de la cadena
			}
		} else {
			result.WriteByte(char)
		}

		i++
	}

	return result.String(), nil
}

// cleaner realiza la limpieza y preparación de la expresión regular
func cleaner(regex string) string {
	var regexWithConcatSymbol strings.Builder

	for i, r := range regex {
		regexWithConcatSymbol.WriteRune(r)

		if i+1 < len(regex) {
			nextRune := rune(regex[i+1])

			if !strings.ContainsRune("(|", r) && !strings.ContainsRune("*+?|)", nextRune) {
				regexWithConcatSymbol.WriteRune('^')
			}
		}
	}
	return regexWithConcatSymbol.String()
}

func shuntingYard(infix string) string {
	precedence := map[rune]int{
		'*': 4, '+': 4, '?': 4, '^': 3, '|': 2, '(': 1,
	}
	fmt.Println(infix)
	infix = cleaner(infix)
	fmt.Println(infix)
	postfix := ""
	stack := []rune{}

	for _, char := range infix {
		if unicode.IsLetter(char) || unicode.IsNumber(char) || char == 'ε' {
			postfix += string(char)
		} else if char == '(' {
			stack = append(stack, char)
		} else if char == ')' {
			for len(stack) > 0 && stack[len(stack)-1] != '(' {
				postfix += string(stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
			if len(stack) > 0 {
				stack = stack[:len(stack)-1] // Pop '('
			}
		} else if prec, ok := precedence[char]; ok {
			for len(stack) > 0 {
				peek := stack[len(stack)-1]
				if precedence[peek] >= prec {
					postfix += string(peek)
					stack = stack[:len(stack)-1]
				} else {
					break
				}
			}
			stack = append(stack, char)
		}
	}

	for len(stack) > 0 {
		postfix += string(stack[len(stack)-1])
		stack = stack[:len(stack)-1]
	}
	fmt.Printf("Shunting Yard %s\n", postfix)
	return postfix
}

func extendedShuntingYard(infix utils.DoublyLinkedList, definitions map[string]utils.DoublyLinkedList) ([]utils.RegexToken, error){
	// infix.PrintForward()
	precedence := map[string][]int{
		"KLEENE": {4,42}, "CATOPERATOR": {3,94}, "OROPERATOR": {2,124}, "OPENPARENTHESES": {1,40},
	}
	posfixExp := []utils.RegexToken{}
	operatorStack := utils.Stack{}

	currentToken := infix.Head
	for currentToken != nil{
		// fmt.Println(posfixExp) 
		// fmt.Printf("Top of stack%s\n", operatorStack.Peek())
		if currentToken.Value.(utils.RegexToken).IsOperator != ""{
			operator :=currentToken.Value.(utils.RegexToken).IsOperator
			switch operator{
			case "KLEENE", "CATOPERATOR", "OROPERATOR":
				// fmt.Printf("Case %s\n", operator)
				for operatorStack.Size() > 0 && precedence[operator][0] <= precedence[operatorStack.Peek().(string)][0]{
					// fmt.Print("--------ACA------\n")
					tempOperator := operatorStack.Pop().(string)
					posfixExp = append(posfixExp, utils.RegexToken{
						Value: []rune{rune(precedence[tempOperator][1])},
						IsOperator: tempOperator,
					})
					// fmt.Printf("Top of stack %s\n", operatorStack.Peek())
				}
				operatorStack.Push(operator)
			case "OPENPARENTHESES":
				operatorStack.Push(operator)
			case "CLOSEPARENTHESES":
				for operatorStack.Size() > 0 && operatorStack.Peek() != "OPENPARENTHESES" {
					tempOperator := operatorStack.Pop().(string)
					posfixExp = append(posfixExp, utils.RegexToken{
						Value: []rune{rune(precedence[tempOperator][1])},
						IsOperator: tempOperator,
					})
				}
				operatorStack.Pop()
			default:
				// fmt.Printf("Case default, op not recognized %s\n", operator)
				if operator == "NULL"{
					posfixExp = append(posfixExp, currentToken.Value.(utils.RegexToken))
				}else if definitions[operator].Head == nil {
					return nil, fmt.Errorf("regex parsing error: indent %s not recognized", operator)
				}else{
					nextToken := currentToken.Next
					currentToken.Next = definitions[operator].Head
					currentToken.Next.Prev = currentToken
					nextToken.Prev = definitions[operator].Tail
					nextToken.Prev.Next = nextToken
				}
			}
		}else{
			// fmt.Printf("No ope found\n")
			posfixExp = append(posfixExp, currentToken.Value.(utils.RegexToken))
		}
		currentToken = currentToken.Next
	}
	for operatorStack.Size() > 0{
		tempOperator := operatorStack.Pop().(string)
		posfixExp = append(posfixExp, utils.RegexToken{
			Value: []rune{rune(precedence[tempOperator][1])},
			IsOperator: tempOperator,
		})
	}


	// fmt.Println(posfixExp)
	// fmt.Printf("Top of stack %s\n\n", operatorStack.Peek())

	return posfixExp, nil
}

func InfixToPosfix(regex string) (string, error) {
	validatedRegex, err := validation(regex)
	if err != nil {
		return "", err
	} else {
		convertedRegex, err := convertQuestionMarkAndPlusSign(validatedRegex)
		if err != nil {
			return "", err
		} else {
			return shuntingYard(convertedRegex), nil
		}
	}
}

func ExtendedInfixToPosfix(regex utils.DoublyLinkedList, definitions map[string]utils.DoublyLinkedList) ([]utils.RegexToken, error) {
	return extendedShuntingYard(regex, definitions)	
}

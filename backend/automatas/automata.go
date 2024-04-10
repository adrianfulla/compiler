package automatas

import (
	"fmt"
	"strconv"
	"strings"
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

func escapeRune(char rune) (rune, error){
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

func extendedValidation(regex string) (string, error) {
	if len(regex) == 0 {
		return "", fmt.Errorf("ingrese una expresión regular")
	}
	regexDLinkedList := utils.DoublyLinkedList{}
	
	for _, char := range regex {
		regexDLinkedList.Append(utils.RegexToken{
			Value: []rune{char},
		})
	}
	currentToken := regexDLinkedList.Head
	for currentToken != nil {
		fmt.Println("")
		fmt.Println("Forward")
		regexDLinkedList.PrintForward()
		fmt.Println("Reverse")
		regexDLinkedList.PrintReverse()
		char := currentToken.Value.(utils.RegexToken).Value
		if len(char) == 1{
			char := char[len(char)-1]
			fmt.Printf("Encontrado:%s\n", string(char))
			switch char{
			case '\'':
				fmt.Printf("Case:%s\n", string(char))
				if currentToken.Next == nil{
					return "", fmt.Errorf("regex parsing error, found extra '")
				}
				prevCharToken :=utils.RegexToken{
					Value: []rune{'('},
					IsOperator: "OPENPARENTHESES",
				}
				if currentToken.Prev == nil{
					regexDLinkedList.Head = currentToken.Next
					regexDLinkedList.Prepend(prevCharToken)
					currentToken = regexDLinkedList.Head
				}else{
					newPrevToken := utils.LinkedNode{
						Value: prevCharToken,
						Next: currentToken.Next,
						Prev: currentToken.Prev,
					}
					currentToken.Next.Prev = &newPrevToken
					currentToken.Prev.Next = &newPrevToken
					currentToken = currentToken.Prev.Prev.Next.Next
				}

				// fmt.Printf("Current Token %s \n", currentToken.Value)
				// fmt.Printf("Next Token %s \n", currentToken.Next.Value)
				// fmt.Printf("Next Token %s \n", currentToken.Next.Next.Value)
				
				currentToken = currentToken.Next

				catChar := utils.RegexToken{
					Value: []rune{'^'},
					IsOperator:"CATOPERATOR",
				}
				for currentToken.Next.Value.(utils.RegexToken).Value[0] != '\''{
					newNextToken := utils.LinkedNode{
						Value: catChar,
						Next: currentToken.Next,
						Prev: currentToken,
					}
					currentToken.Next.Prev = &newNextToken
					currentToken.Next = &newNextToken
					currentToken = newNextToken.Next
				}
				// fmt.Printf("Current Token %s \n", currentToken.Value)
				// fmt.Printf("Next Token %s \n", currentToken.Next.Value)
				
				nextChar := utils.RegexToken{
					Value: []rune{')'},
					IsOperator: "CLOSEPARENTHESES",
				}
				if currentToken.Next.Next == nil{
					// fmt.Println("Encontro null")
					// fmt.Printf("Current Token %s \n", currentToken.Value)
					// fmt.Printf("Next Token %s \n", currentToken.Next.Value)
					regexDLinkedList.Tail = currentToken
					regexDLinkedList.Append(nextChar)
					
				}else{
					newNextToken := utils.LinkedNode{
						Value: nextChar,
						Next: currentToken.Next.Next,
						Prev: currentToken,
					}
					currentToken.Next.Next.Prev = &newNextToken
					currentToken.Next = &newNextToken
				}
				currentToken = currentToken.Next
			case '"':
				fmt.Printf("Case:%s\n", string(char))
				if currentToken.Next == nil{
					return "", fmt.Errorf("regex parsing error, found extra '")
				}
				prevCharToken :=utils.RegexToken{
					Value: []rune{'('},
					IsOperator: "OPENPARENTHESES",
				}
				if currentToken.Prev == nil{
					regexDLinkedList.Head = currentToken.Next
					regexDLinkedList.Prepend(prevCharToken)
					currentToken = regexDLinkedList.Head
				}else{
					newPrevToken := utils.LinkedNode{
						Value: prevCharToken,
						Next: currentToken.Next,
						Prev: currentToken.Prev,
					}
					currentToken.Next.Prev = &newPrevToken
					currentToken.Prev.Next = &newPrevToken
					currentToken = currentToken.Prev.Prev.Next.Next
				}

				// fmt.Printf("Current Token %s \n", currentToken.Value)
				// fmt.Printf("Next Token %s \n", currentToken.Next.Value)
				// fmt.Printf("Next Token %s \n", currentToken.Next.Next.Value)
				
				currentToken = currentToken.Next

				catChar := utils.RegexToken{
					Value: []rune{'^'},
					IsOperator:"CATOPERATOR",
				}
				for currentToken.Next.Value.(utils.RegexToken).Value[0] != '"'{
					newNextToken := utils.LinkedNode{
						Value: catChar,
						Next: currentToken.Next,
						Prev: currentToken,
					}
					currentToken.Next.Prev = &newNextToken
					currentToken.Next = &newNextToken
					currentToken = newNextToken.Next
				}
				// fmt.Printf("Current Token %s \n", currentToken.Value)
				// fmt.Printf("Next Token %s \n", currentToken.Next.Value)
				
				nextChar := utils.RegexToken{
					Value: []rune{')'},
					IsOperator: "CLOSEPARENTHESES",
				}
				if currentToken.Next.Next == nil{
					// fmt.Println("Encontro null")
					// fmt.Printf("Current Token %s \n", currentToken.Value)
					// fmt.Printf("Next Token %s \n", currentToken.Next.Value)
					regexDLinkedList.Tail = currentToken
					regexDLinkedList.Append(nextChar)
					
				}else{
					newNextToken := utils.LinkedNode{
						Value: nextChar,
						Next: currentToken.Next.Next,
						Prev: currentToken,
					}
					currentToken.Next.Next.Prev = &newNextToken
					currentToken.Next = &newNextToken
				}
				currentToken = currentToken.Next
			case '_':
				fmt.Printf("Case:%s\n", string(char))
				
				prevChar :=utils.RegexToken{
					Value: []rune{'('},
					IsOperator: "OPENPARENTHESES",
				}
				nextChar := utils.RegexToken{
					Value: []rune{')'},
					IsOperator: "CLOSEPARENTHESES",
				}
				if currentToken.Prev == nil{
					regexDLinkedList.Prepend(prevChar)
				}else{
					newPrevToken := utils.LinkedNode{
						Value: prevChar,
						Next: currentToken,
						Prev: currentToken.Prev,
					}
					currentToken.Prev.Next = &newPrevToken
				}

				orChar := utils.RegexToken{
					Value: []rune{'|'},
					IsOperator: "OROPERATOR",
				}
				nextToken := currentToken.Next
				fmt.Printf("Next Token %s \n", currentToken.Next.Value)
				for i := 0; i <= 254; i++{
					// fmt.Printf("Including %s\n", rune(i))
					currentToken.Value = utils.RegexToken{
						Value: []rune{rune(i)},
					}
					
					currentToken.Next = &utils.LinkedNode{
						Value: orChar,
						Next: &utils.LinkedNode{},
						Prev: currentToken,
					}
					currentToken.Next.Next.Prev = currentToken
					currentToken.Next.Next.Next = nextToken

					currentToken = currentToken.Next.Next
				}
				// fmt.Printf("Prev Token %s \n", currentToken.Prev.Value)
				// fmt.Printf("Current Token %s \n", currentToken.Value)
				// fmt.Printf("Next Token %s \n", currentToken.Next.Value)
				
				currentToken.Value = utils.RegexToken{
					Value: []rune{rune(255)},
				}
				if currentToken.Next == nil{
					regexDLinkedList.Append(nextChar)
				}else{
					newNextToken := utils.LinkedNode{
						Value: nextChar,
						Next: currentToken.Next,
						Prev: currentToken,
					}
					currentToken.Next.Prev = &newNextToken
					currentToken.Next = &newNextToken
				}
			case '[':
			// 	fmt.Printf("Case:%s\n", string(char))
			// 	if currentToken.Next == nil{
			// 		return "", fmt.Errorf("regex parsing error, found extra '")
			// 	}
			// 	prevCharToken :=utils.RegexToken{
			// 		Value: []rune{'('},
			// 		IsOperator: "OPENPARENTHESES",
			// 	}
			// 	if currentToken.Prev == nil{
			// 		regexDLinkedList.Head = currentToken.Next
			// 		regexDLinkedList.Prepend(prevCharToken)
			// 		currentToken = regexDLinkedList.Head
			// 	}else{
			// 		newPrevToken := utils.LinkedNode{
			// 			Value: prevCharToken,
			// 			Next: currentToken.Next,
			// 			Prev: currentToken.Prev,
			// 		}
			// 		currentToken.Next.Prev = &newPrevToken
			// 		currentToken.Prev.Next = &newPrevToken
			// 		currentToken = currentToken.Prev.Prev.Next.Next
			// 	}

			// 	// fmt.Printf("Current Token %s \n", currentToken.Value)
			// 	// fmt.Printf("Next Token %s \n", currentToken.Next.Value)
			// 	// fmt.Printf("Next Token %s \n", currentToken.Next.Next.Value)
				
			// 	currentToken = currentToken.Next

			// 	catChar := utils.RegexToken{
			// 		Value: []rune{'^'},
			// 		IsOperator:"CATOPERATOR",
			// 	}
			// 	for currentToken.Next.Value.(utils.RegexToken).Value[0] != '"'{
			// 		newNextToken := utils.LinkedNode{
			// 			Value: catChar,
			// 			Next: currentToken.Next,
			// 			Prev: currentToken,
			// 		}
			// 		currentToken.Next.Prev = &newNextToken
			// 		currentToken.Next = &newNextToken
			// 		currentToken = newNextToken.Next
			// 	}
			// 	// fmt.Printf("Current Token %s \n", currentToken.Value)
			// 	// fmt.Printf("Next Token %s \n", currentToken.Next.Value)
				
			// 	nextChar := utils.RegexToken{
			// 		Value: []rune{')'},
			// 		IsOperator: "CLOSEPARENTHESES",
			// 	}
			// 	if currentToken.Next.Next == nil{
			// 		// fmt.Println("Encontro null")
			// 		// fmt.Printf("Current Token %s \n", currentToken.Value)
			// 		// fmt.Printf("Next Token %s \n", currentToken.Next.Value)
			// 		regexDLinkedList.Tail = currentToken
			// 		regexDLinkedList.Append(nextChar)
					
			// 	}else{
			// 		newNextToken := utils.LinkedNode{
			// 			Value: nextChar,
			// 			Next: currentToken.Next.Next,
			// 			Prev: currentToken,
			// 		}
			// 		currentToken.Next.Next.Prev = &newNextToken
			// 		currentToken.Next = &newNextToken
			// 	}
			// 	currentToken = currentToken.Next
			case ']':
				return "", fmt.Errorf("unexpected closing bracket found at regex")
			case '(':
			case ')':
			case '^':
			case '#':
			case '*':
			case '+':
			case '?':
			case '|':
			
			default:
				fmt.Printf("Default Case:%s\n", string(char))
			// 	if currentToken.Next == nil{
			// 		return "", fmt.Errorf("regex parsing error, found extra '")
			// 	}
			// 	nextChar := currentToken.Next.Value.(utils.RegexToken).Value[0]
			// 	currentToken.Value = utils.RegexToken{
			// 		Value: []rune{},
			// 	}
			// 	for nextChar != '"'{
			// 		if currentToken.Next != nil{
			// 			if nextChar == '\\'{
			// 				temp, err := escapeRune(currentToken.Next.Next.Value.(utils.RegexToken).Value[0])
			// 				if err != nil{
			// 					if currentToken.Next.Next.Value.(utils.RegexToken).Value[0] == '"'{
			// 						temp = '\''
			// 					}else{
			// 						return "", fmt.Errorf("invalid escape sequence in regex")
			// 					}
			// 				}
			// 				nextChar = temp
			// 				currentToken.Next = currentToken.Next.Next
			// 			}
			// 			currentToken.Value = utils.RegexToken{
			// 				Value: append(currentToken.Value.(utils.RegexToken).Value, nextChar),
			// 			}
			// 			currentToken.Next = currentToken.Next.Next
			// 			nextChar = currentToken.Next.Value.(utils.RegexToken).Value[0]
			// 		}else{
			// 			break
			// 		}
			// 	}
			// 	if currentToken.Prev != nil{
			// 		currentToken.Prev.Next = currentToken
			// 	}else{
			// 		regexDLinkedList.Head = currentToken
			// 	}
			// 	currentToken.Next = currentToken.Next.Next
			// }
		// }
		
			}
			currentToken = currentToken.Next
		}
	}
	fmt.Println("Forward")
	regexDLinkedList.PrintForward()
	fmt.Println("Reverse")
	regexDLinkedList.PrintReverse()
	return regex, nil
}


func expandBrackets(limits []rune, isNot bool) (utils.Stack, error) {
	result := utils.Stack{}
	start := limits[1]
	end := limits[0]
	if isNot {
		if !strings.ContainsRune("Aa0", start) {
			runeType := utils.RuneType(start)
			var val rune
			switch runeType {
			case "upper":
				val = 'A'
			case "lower":
				val = 'a'
			case "number":
				val = '0'
			case "other":
				return result, fmt.Errorf("character in character set not allowed for start: %s", string(start))
			}
			tStart := val
			tEnd := start - 1
			res, err := expandBrackets([]rune{tEnd, tStart}, false)
			if err != nil {
				return result, err
			}
			for res.Size() > 0 {
				result.Push(res.Pop())
			}
		}
		if !strings.ContainsRune("Zz9", end) {
			runeType := utils.RuneType(start)
			var val rune
			switch runeType {
			case "upper":
				val = 'Z'
			case "lower":
				val = 'z'
			case "number":
				val = '9'
			case "other":
				return result, fmt.Errorf("character in character set not allowed for end: %s", string(start))
			}
			tStart := end + 1
			tEnd := val
			res, err := expandBrackets([]rune{tEnd, tStart}, false)
			if err != nil {
				return result, err
			}
			for res.Size() > 0 {
				result.Push(res.Pop())
			}
		}
		return result, nil
	}
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
	return postfix
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

func ExtendedInfixToPosfix(regex string) (string, error) {
	validatedRegex, err := extendedValidation(regex)
	if err != nil {
		fmt.Print(err)
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

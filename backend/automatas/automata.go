package automatas

import (
	"fmt"
	"math"
	"strings"
	"unicode"

	"github.com/adrianfulla/compiler/backend/utils"
)

type valChar struct {
	value     rune
	isEscaped bool
}
func Difference(a, b []valChar) ([]rune) {
    // Crear mapas para guardar los elementos de cada array
    mapA := make(map[rune]bool)
    mapB := make(map[rune]bool)
	diffAB := []rune{}


    // Llenar el mapa para el primer array
    for _, item := range a {
		
        mapA[item.value] = true
    }

    // Llenar el mapa para el segundo array
    for _, item := range b {
        mapB[item.value] = true
    }

    // Encontrar elementos que están en A pero no en B
    for _, item := range a {
        if !mapB[item.value] {
            diffAB = append(diffAB, item.value)
        }
    }

    // Encontrar elementos que están en B pero no en A
    for _, item := range b {
		// fmt.Printf("in b: %s\n", string(item.value))
        if !mapA[item.value] {
			// fmt.Printf("in b: %s\n", string(item.value))
            diffAB = append(diffAB, item.value)
        }
    }

    return diffAB
}

func appendChar(exp strings.Builder, char valChar) strings.Builder {
	temp := exp.String()
	exp.Reset()
	// fmt.Printf("Writing %s\n", string(char.value))
	exp.WriteRune('\'')
	if char.isEscaped {
		exp.WriteRune('\\')
	}
	exp.WriteRune(char.value)
	exp.WriteRune('\'')
	exp.WriteString(temp)
	return exp
}
func appendCharOR(exp strings.Builder, char valChar) strings.Builder {
	temp := exp.String()
	exp.Reset()
	// fmt.Printf("Writing %s\n", string(char.value))
	exp.WriteRune('\'')
	if char.isEscaped {
		exp.WriteRune('\\')
	}
	exp.WriteRune(char.value)
	exp.WriteRune('\'')
	exp.WriteRune('|')
	exp.WriteString(temp)
	return exp
}

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

func extendedValidation(regex string) (string, error) {
	if len(regex) == 0 {
		return "", fmt.Errorf("ingrese una expresión regular")
	}
	var expandedRegex strings.Builder
	inApostrophe := false
	inDoubleApostrophe := false
	isEscaped := false
	inCharSet := false
	inParentheses := false
	prevCharacterSet := []valChar{}
	isDiffofCharSets := false
	isNotCharacterSet := false
	stack := utils.Stack{}
	for _, char := range regex {
		fmt.Printf("expanded: %s\n", expandedRegex.String())
		fmt.Printf("Char: %s,isEsc: %t, isApo: %t, isDoubleApo: %t, inCharSet: %t\n", string(char), isEscaped, inApostrophe, inDoubleApostrophe, inCharSet)
		switch char {
		case '\'':
			if !inCharSet {
				if inApostrophe && !isEscaped {
					if stack.Size() > 1 {
						return "", fmt.Errorf("invalid character in single apostrophe")
					}
					for stack.Size() > 0 {
						expandedRegex = appendChar(expandedRegex, stack.Pop().(valChar))
					}
					inApostrophe = false
				} else if isEscaped || inDoubleApostrophe {
					stack.Push(valChar{
						value:     char,
						isEscaped: isEscaped,
					})
					isEscaped = false
				} else if !inDoubleApostrophe {
					inApostrophe = true
				} else {
					return "", fmt.Errorf("invalid character: %s", string(char))
				}
			}
		case '"':
			if inDoubleApostrophe && !isEscaped {
				for stack.Size() > 0 {
					expandedRegex = appendChar(expandedRegex, stack.Pop().(valChar))
				}
				inDoubleApostrophe = false
			} else if inApostrophe || isEscaped {
				stack.Push(valChar{
					value:     char,
					isEscaped: isEscaped,
				})
				isEscaped = false
			} else if !inApostrophe {
				inDoubleApostrophe = true
			} else {
				return "", fmt.Errorf("invalid character: %s", string(char))
			}
		case '\\':
			if !inApostrophe && !inDoubleApostrophe {
				return "", fmt.Errorf("invalid character escape")
			}
			if isEscaped {
				stack.Push(valChar{
					value:     char,
					isEscaped: isEscaped,
				})
				isEscaped = false
			} else {
				isEscaped = true
			}
		case '[':
			if inApostrophe || inDoubleApostrophe {
				stack.Push(valChar{
					value:     char,
					isEscaped: isEscaped,
				})
				isEscaped = false
			} else if inCharSet {
				return "", fmt.Errorf("invalid character: %s", string(char))
			} else {
				inCharSet = true
			}
		case ']':
			if !inCharSet && !(inApostrophe || inDoubleApostrophe) {
				return "", fmt.Errorf("invalid character: %s", string(char))
			} else if  inApostrophe ||inDoubleApostrophe{
				stack.Push(valChar{
					value:     char,
					isEscaped: isEscaped,
				})
				isEscaped = false
			} else if !inDoubleApostrophe {
				if isDiffofCharSets{
					fmt.Printf("diff of charset\n")
					newCharacterSet := []valChar{}
					for stack.Size() > 0{
						if len(newCharacterSet) < 2 {
							newCharacterSet = append(newCharacterSet, stack.Pop().(valChar))
						}
						if len(newCharacterSet) == 2 {
							newCharSet, _ := expandBrackets(newCharacterSet, isNotCharacterSet)
							newCharacterSet = []valChar{}
							for newCharSet.Size() > 0{
								newCharacterSet = append(newCharacterSet, newCharSet.Pop().(valChar))
							}
						}
					}
					oldCharSet, _ := expandBrackets(prevCharacterSet, false)
					oldCharacterSet := []valChar{}
					for oldCharSet.Size() > 0{
						oldCharacterSet = append(oldCharacterSet, oldCharSet.Pop().(valChar))
					}
					diff := Difference(oldCharacterSet, newCharacterSet)
					tExp := strings.Builder{}
					for x, i := range diff{
						if x != 0{
							tExp = appendCharOR(tExp, valChar{
								value: i,
								isEscaped: false,
							})
						}else{
							tExp = appendChar(tExp, valChar{
								value: i,
								isEscaped: false,
							})
						}
						
					}
					expandedRegex = tExp
					
				}else if stack.Size() == 1 {
					expandedRegex = appendChar(expandedRegex, stack.Pop().(valChar))
				} else if math.Mod(float64(stack.Size()), 2) == 0 {
					characterSet := []valChar{}
					fmt.Printf("Stack size %d\n", stack.Size())
					for stack.Size() > 0 {
						if len(characterSet) < 2 {
							characterSet = append(characterSet, stack.Pop().(valChar))
						}
						if len(characterSet) == 2 {
							tCharSet, err := expandBrackets(characterSet, isNotCharacterSet)
							prevCharacterSet = append(prevCharacterSet, characterSet...)
							holdString := strings.Builder{}
							if err != nil {
								fmt.Printf("Error\n")
								return "", fmt.Errorf("invalid character set")
							}
							holdString = appendChar(holdString, tCharSet.Pop().(valChar))
							for tCharSet.Size() > 0 {
								holdString = appendCharOR(holdString, tCharSet.Pop().(valChar))
							}
							t := strings.Builder{}
							t.WriteString(expandedRegex.String())
							t.WriteString(holdString.String())
							expandedRegex = t 
							characterSet = []valChar{}
						}
					}
				}
				
				inCharSet = false
			}
		case '-':
			if !inCharSet {
				stack.Push(valChar{
					value:     char,
					isEscaped: isEscaped,
				})
				isEscaped = false
			}
		case '^':
			if !inCharSet && (!inApostrophe || !inDoubleApostrophe) {
				return "", fmt.Errorf("invalid character: %s", string(char))
			} else if inCharSet {
				if !isNotCharacterSet {
					isNotCharacterSet = true
				} else {
					return "", fmt.Errorf("invalid character: %s", string(char))
				}
			} else {
				stack.Push(valChar{
					value:     char,
					isEscaped: isEscaped,
				})
				if isEscaped {
					isEscaped = false
				}
			}
		case '#':
			if inApostrophe|| inDoubleApostrophe  {
				stack.Push(valChar{
					value:     char,
					isEscaped: isEscaped,
				})
				isEscaped = false
			}else{
				isDiffofCharSets = true
			}
		case '(':
			if !inParentheses && (!inApostrophe || !inDoubleApostrophe){
				inParentheses = true
				temp := expandedRegex.String()
				expandedRegex = strings.Builder{}
				expandedRegex.WriteString(temp+string(char))
			}else{
				stack.Push(valChar{
					value:     char,
					isEscaped: isEscaped,
				})
				if isEscaped {
					isEscaped = false
				}
			}
		case ')':
			if !inParentheses && !(inApostrophe || inDoubleApostrophe) {
				return "", fmt.Errorf("invalid character: %s", string(char))
			} else if  inApostrophe ||inDoubleApostrophe{
				stack.Push(valChar{
					value:     char,
					isEscaped: isEscaped,
				})
				isEscaped = false
			} else if !inDoubleApostrophe {
				count := 1
				for inParentheses{
					temp := strings.Builder{}
					str := expandedRegex.String()
					val := str[len(str)-count:len(str)-1]
					fmt.Printf("exp in par %s\n",val)
					temp.WriteString(str)
					if val == "(" {
						inParentheses = false
						fmt.Printf("Inside par: %s\n", temp.String())
					}else{
						temp.WriteString(val)
						count++
					}
					// fmt.Printf("string in parentheses%s\n", )
				}
				temp := expandedRegex.String()
				expandedRegex = strings.Builder{}
				expandedRegex.WriteString(temp+string(char))
			}
		default:
			if inApostrophe || inDoubleApostrophe || inCharSet{
				stack.Push(valChar{
					value:     char,
					isEscaped: isEscaped,
				})
				if isEscaped {
					isEscaped = false
				}
			}else {
				temp := expandedRegex.String()
				expandedRegex = strings.Builder{}
				expandedRegex.WriteString(temp+string(char))
			}
		}
	}
	fmt.Printf("Finished Expanded: %s\n", expandedRegex.String())
	// stack := []rune{}
	// stackOp := []rune{}
	// inApostrophe := false
	// for _, char := range regex {
	// 	fmt.Printf("char in Validation: %s\n", string(char))
	// 	switch char {
	// 	case '(':
	// 		if !inApostrophe{
	// 			stack = append(stack, char)
	// 		}
	// 	case ')':
	// 		if !inApostrophe{
	// 			if len(stack) == 0 {
	// 				return "", fmt.Errorf("paréntesis no balanceados en la expresión regular")
	// 			}
	// 			stack = stack[:len(stack)-1] // Simula el pop
	// 		}

	// 	case '\'':
	// 		inApostrophe = !inApostrophe
	// 	}
	// 	if !inApostrophe{
	// 		stackOp = append(stackOp, char)
	// 	}
	// }

	// if len(stack) > 0 {
	// 	return "", fmt.Errorf("paréntesis no balanceados en la expresión regular")
	// }

	// operators := "*+?|"
	// for i := 0; i < len(regex)-1; i++ {
	// 	currentChar := regex[i]
	// 	nextChar := regex[i+1]

	// 	if strings.ContainsRune(operators, rune(currentChar)) && strings.ContainsRune(operators, rune(nextChar)) && nextChar != '|' {
	// 		return "", fmt.Errorf("sintaxis incorrecta de operadores en la expresión regular")
	// 	}
	// }
	return regex, nil
}

func expandBrackets(limits []valChar, isNot bool) (utils.Stack, error) {
	result := utils.Stack{}
	start := limits[1].value
	end := limits[0].value
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
			tStart := valChar{
				value:     val,
				isEscaped: false,
			}
			tEnd := valChar{
				value:     start - 1,
				isEscaped: false,
			}
			res, err := expandBrackets([]valChar{tEnd, tStart}, false)
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
			tStart := valChar{
				value:     end + 1,
				isEscaped: false,
			}
			tEnd := valChar{
				value:     val,
				isEscaped: false,
			}
			res, err := expandBrackets([]valChar{tEnd, tStart}, false)
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
		result.Push(valChar{
			value:     c,
			isEscaped: false,
		})
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

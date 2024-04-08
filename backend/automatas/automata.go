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
	var expandedRegex strings.Builder
	inApostrophe := false
	inDoubleApostrophe := false
	isEscaped := false
	inCharSet := false
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
				}
				if stack.Size() == 1 {
					expandedRegex = appendChar(expandedRegex, stack.Pop().(valChar))
				} else if math.Mod(float64(stack.Size()), 2) == 0 {
					characterSet := []valChar{}
					fmt.Printf("Stack size %d\n", stack.Size())
					for stack.Size() > 0 {
						fmt.Printf("CharSet size %d\n", len(characterSet))
						if len(characterSet) < 2 {
							characterSet = append(characterSet, stack.Pop().(valChar))
						}
						if len(characterSet) == 2 {
							tCharSet, err := expandBrackets(characterSet, isNotCharacterSet)
							prevCharacterSet = append(prevCharacterSet, characterSet...)
							if err != nil {
								fmt.Printf("Error\n")
								return "", fmt.Errorf("invalid character set")
							}
							for tCharSet.Size() > 0 {
								expandedRegex = appendCharOR(expandedRegex, tCharSet.Pop().(valChar))
							}
							characterSet = []valChar{}
						}
					}
				}
			} else {
				
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
			if !inCharSet {
				stack.Push(valChar{
					value:     char,
					isEscaped: isEscaped,
				})
				isEscaped = false
			}else{
				isDiffofCharSets = true
			}
		default:
			stack.Push(valChar{
				value:     char,
				isEscaped: isEscaped,
			})
			if isEscaped {
				isEscaped = false
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

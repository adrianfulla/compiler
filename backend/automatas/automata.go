package automatas

import (
	"fmt"
	"strings"
	"unicode"
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
				result.WriteString( string(prevChar) + "*")
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

			// if (unicode.IsLetter(r) || unicode.IsNumber(r) || r == 'ε' || r == ')') &&
			// 	(unicode.IsLetter(nextRune) || unicode.IsNumber(nextRune) || nextRune == '(' || nextRune == 'ε') {
			// 	regexWithConcatSymbol.WriteRune('^')
			// }
			if (!strings.ContainsRune("(|", r) && !strings.ContainsRune("*+?|)", nextRune)) {
				regexWithConcatSymbol.WriteRune('^')
			}
		}
	}
	return regexWithConcatSymbol.String()
}



// shouldConcatenate decide si se debe insertar un operador de concatenación entre dos caracteres
func shouldConcatenate(prevRune, currentRune, nextRune rune) bool {
	isCurrentAlphaNumOrEpsilon := unicode.IsLetter(currentRune) || unicode.IsNumber(currentRune) || currentRune == 'ε'
	isNextAlphaNumOrEpsilon := unicode.IsLetter(nextRune) || unicode.IsNumber(nextRune) || nextRune == 'ε'
	isNextOpenParenthesis := nextRune == '('
	isPrevOROperator := prevRune == '|'
	isPrevCloseParenthesis := prevRune == ')'

	// Agrega la concatenación solo si el siguiente caracter es alfanumérico, un épsilon, o un paréntesis abierto,
	// y el actual es alfanumérico o un épsilon, pero no si el siguiente es un operador (excluyendo el paréntesis abierto).
	return isCurrentAlphaNumOrEpsilon && (isNextAlphaNumOrEpsilon || isNextOpenParenthesis) && (!isPrevOROperator ||isPrevCloseParenthesis )
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

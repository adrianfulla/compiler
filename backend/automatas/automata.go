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

	if strings.Contains(regex, "|") {
		parts := strings.Split(regex, "|")
		for _, part := range parts {
			if strings.ContainsAny(part, "()*+?.") {
				return "", fmt.Errorf("operador OR mal utilizado en la expresión regular")
			}
		}
	}
	return regex, nil
}

// cleaner realiza la limpieza y preparación de la expresión regular
func cleaner(regex string) string {
	// Reemplaza los símbolos especiales para simplificar la conversión a posfijo
	firstCleaned := regex
	for i := 1; i < len(firstCleaned); i++ {
		if firstCleaned[i] == '+' && (i+1 < len(firstCleaned) && firstCleaned[i+1] != '*') {
			temp := string(firstCleaned[i-1])
			j := i - 1
			for j > 0 && firstCleaned[j-1] == firstCleaned[i-1] {
				temp += string(firstCleaned[j-1])
				j--
			}
			firstCleaned = strings.Replace(firstCleaned, temp+"+", temp+temp+"*", 1)
		}
	}

	// Agrega el símbolo de concatenación explícito
	regexWithConcatSymbol := ""
	for i := 0; i < len(firstCleaned); i++ {
		regexWithConcatSymbol += string(firstCleaned[i])
		if i+1 < len(firstCleaned) {
			nextChar := firstCleaned[i+1]
			if firstCleaned[i] != '(' && nextChar != ')' && !strings.ContainsRune("*+?|", rune(nextChar)) && firstCleaned[i] != '|' {
				regexWithConcatSymbol += "^"
			}
		}
	}

	return regexWithConcatSymbol
}


func shuntingYard(infix string) string {
	precedence := map[rune]int{
		'*': 4, '+': 4, '?': 4, '^': 3, '|': 2, '(': 1,
	}

	infix = cleaner(infix)
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
				stack = stack[:len(stack)-1]
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
	if err!= nil {
        return "", err
    } else {
		return shuntingYard(validatedRegex), nil
	}
}
package parser

import (
	// "fmt"

	// "fmt"

	"fmt"
	"sort"
	"encoding/json"
	"strings"
	"sync"
	"strconv"

	"github.com/adrianfulla/compiler/backend/automatas"
	"github.com/adrianfulla/compiler/backend/lexer"
	"github.com/adrianfulla/compiler/backend/utils"
)

type Parser struct {
	file          string
	scanner       *lexer.Scanner
	afdStack      map[string]automatas.DAfdJson
	Tokens        []utils.ParseToken      `json:"tokens"`
	IgnoredTokens []utils.ParseToken      `json:"productions"`
	Productions   []*utils.ProductionToken `json:"ignored_tokens"`
	SLR *automatas.SLR `json:"slr"`
}

func LexYaparFile(fileYapar string, lexScanner *lexer.Scanner) (*Parser, error) {
	// fmt.Print(fileYapar)
	lex := &Parser{
		file:     fileYapar,
		afdStack: make(map[string]automatas.DAfdJson),
		scanner:  lexScanner,
	}

	definitions := map[string]string{}

	definitions["COMMENT"] = "'/* '['A'-'Z''a'-'z''0'-'9'\" .\"]*' */'"
	definitions["LOWERCASE"] = "['a'-'z']+"
	definitions["UPPERCASE"] = "'I'['A'-'H''J'-'Z']+|['A'-'H''J'-'Z']['A'-'Z']*"
	definitions["TOKEN"] = "'%'\"token\""
	definitions["IGNOREFLAG"] = "'IGNORE '"
	definitions["TWODOTS"] = "\":\""
	definitions["SEMICOLON"] = "';'"
	definitions["OR"] = "'|'"
	definitions["SPLITTER"] = "'%''%'"
	definitions["SPACE"] = "[' ''\\t']"
	definitions["NEWLINE"] = "['\\n']"

	validatedDefinitions := map[string]*utils.DoublyLinkedList{}

	for token, def := range definitions {
		validated, err := automatas.ExtendedValidation(def)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		validatedDefinitions[token] = validated
	}
	posfixDefinitions := map[string][]utils.RegexToken{}
	for token, def := range validatedDefinitions {
		posfix, err := automatas.ExtendedInfixToPosfix(def, validatedDefinitions)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		posfixDefinitions[token] = posfix
	}
	// fmt.Print(posfixDefinitions)
	for token, posfix := range posfixDefinitions {
		// fmt.Print(token)
		afd := automatas.ExtendedNewDirectAfd(posfix)
		afdJson := afd.ToJson()
		// fmt.Println(afdJson)
		lex.afdStack[token] = *afdJson
	}
	// fmt.Print("ACA")
	Scanner, err := lex.parseFile()
	if err != nil {
		return nil, err
	}
	return Scanner, nil
}

func (lex *Parser) parseFile() (*Parser, error) {
	ch := make(chan map[string]utils.Stack, len(lex.afdStack))
	var wg sync.WaitGroup
	// fmt.Print("ACA")
	for index := range lex.afdStack {
		wg.Add(1)
		go func(el string) {
			defer wg.Done()
			lex.searchYalex(el, ch)
		}(index)
	}
	wg.Wait()
	close(ch)
	tokensFound := []*automatas.AcceptedExp{}
	for maps := range ch {
		for index, result := range maps {
			// fmt.Println(index)
			for result.Size() > 0 {
				res := result.Pop().(*automatas.AcceptedExp)
				res.Token = index
				// fmt.Println(res)
				tokensFound = AddOrUpdateExp(res, tokensFound)
			}
		}
	}
	tokensFound = SortTokens(tokensFound)
	tokensFoundStack := utils.Stack{}
	tokensFoundStack.Push(&automatas.AcceptedExp{
		Token: "EOF",
	})
	for _, tokenFound := range tokensFound {
		tokensFoundStack.Push(tokenFound)
	}

	passedSplitter := false
	token := tokensFoundStack.Pop().(*automatas.AcceptedExp)
	for tokensFoundStack.Size() > 0 {
		// fmt.Printf("Case %s %s\n", token.Token, token.Value)
		switch token.Token {
		case "SPLITTER":
			if !passedSplitter {
				passedSplitter = true
			} else {
				return nil, fmt.Errorf("unkown splitter found")
			}

		case "LOWERCASE":
			// fmt.Println("lowercase found")
			if !passedSplitter {
				return nil, fmt.Errorf("error parsing yapar: invalid token definition")
			}
			newProduction := &utils.ProductionToken{
				Head: token.Value,
				Body: [][]string{},
			}
			token = tokensFoundStack.Pop().(*automatas.AcceptedExp)
			if token.Token == "TWODOTS" {
				token = tokensFoundStack.Pop().(*automatas.AcceptedExp)
				body := []string{}
				for token.Token != "SEMICOLON" {
					// fmt.Printf("Case in prod %s %s\n", token.Token, token.Value)
					if token.Token == "OR" {
						if len(body) == 0 {
							return nil, fmt.Errorf("error parsing yapar: invalid OR")
						}
						// fmt.Print(body)
						newProduction.Body = append(newProduction.Body, body)
						token = tokensFoundStack.Pop().(*automatas.AcceptedExp)
						body = []string{}
					}
					body = appendProduction(token, body)
					token = tokensFoundStack.Pop().(*automatas.AcceptedExp)
				}
				newProduction.Body = append(newProduction.Body, body)
			}
			lex.Productions = append(lex.Productions, newProduction)

		case "UPPERCASE":
			// fmt.Println("uppercase found")
			return nil, fmt.Errorf("error parsing yapar: invalid token definition")

		case "TOKEN":
			// fmt.Println("token found")
			if passedSplitter {
				fmt.Print("error parsing yapar: invalid token definition after splitter")
				return nil, fmt.Errorf("error parsing yapar: invalid token definition after splitter")
			}
			token = tokensFoundStack.Pop().(*automatas.AcceptedExp)
			for token.Token != "NEWLINE" {
				// fmt.Printf("Token: %s %s \n", token.Token, token.Value)
				if token.Token == "UPPERCASE" {
					lex.Tokens = append(lex.Tokens, utils.ParseToken{
						Token: token.Value,
					})
				} else {
					if token.Token != "SPACE" {
						return nil, fmt.Errorf("error parsing yapar: invalid token definition")
					}
				}
				token = tokensFoundStack.Pop().(*automatas.AcceptedExp)
			}

		case "IGNOREFLAG":
			// fmt.Println("ignoreflag found")
			token = tokensFoundStack.Pop().(*automatas.AcceptedExp)
			lex.IgnoredTokens = append(lex.IgnoredTokens, utils.ParseToken{
				Token: token.Value,
			})

		case "TWODOTS":
			// fmt.Println("twodots found")
			return nil, fmt.Errorf("error parsing yapar: invalid token definition")

		case "SEMICOLON":
			// fmt.Println("semicolon found")
			return nil, fmt.Errorf("error parsing yapar: invalid token definition")

		case "OR":
			// fmt.Println("lowercase found")
			return nil, fmt.Errorf("error parsing yapar: invalid token definition")

		case "NEWLINE":
			// fmt.Println("newline found")

		default:
			// fmt.Printf("default\n")
		}
		if !tokensFoundStack.IsEmpty() {
			token = tokensFoundStack.Pop().(*automatas.AcceptedExp)
		}
	}

	elementMap := make(map[string]bool)
	for _, item := range lex.Tokens {
		elementMap[item.Token] = true
	}

	// Verificar cada elemento de arr1 en el mapa.
	for _, item := range lex.scanner.Tokens {
		if _, found := elementMap[item.Token]; !found {
			// Si un elemento de arr1 no está en arr2, retornar un error.
			return nil, fmt.Errorf("token in yapar not in yalex")
		}
	}

	// fmt.Print(newScanner.Productions)
	return lex, nil
}

func SortTokens(tokens []*automatas.AcceptedExp) []*automatas.AcceptedExp {
	sort.Slice(tokens, func(i, j int) bool {
		return tokens[i].Start > tokens[j].Start
	})
	return tokens
}

func appendProduction(token *automatas.AcceptedExp, body []string) []string {
	if token.Token != "SPACE" && token.Token != "NEWLINE" && token.Token != "COMMENT" {
		body = append(body, token.Value)
	}
	return body
}

func AddOrUpdateExp(newExp *automatas.AcceptedExp, currentExps []*automatas.AcceptedExp) []*automatas.AcceptedExp {
	temp := []*automatas.AcceptedExp{}
	for _, exp := range currentExps {
		if newExp.Start <= exp.Start && newExp.End >= exp.End {
			continue
		}
		if exp.Start <= newExp.Start && exp.End >= newExp.End {
			return currentExps
		}
		temp = append(temp, exp)
	}
	temp = append(temp, newExp)
	currentExps = temp
	return currentExps
}

func (lex *Parser) searchYalex(index string, ch chan<- map[string]utils.Stack) {
	resultado := make(map[string]utils.Stack)
	resultado[index] = automatas.ExtendedSimulateAfd(lex.file, lex.afdStack[index])
	ch <- resultado
}


func (parser *Parser) PrintParser(){
	fmt.Println("Accepted Tokens:")
	fmt.Println(parser.Tokens)
	fmt.Println("Ignored Tokens:")
	fmt.Println(parser.IgnoredTokens)
	fmt.Print("Productions: [\n")
	for _, prod := range parser.Productions{
		fmt.Printf("\tHead: %s, Body %s\n",
		prod.Head, prod.Body)
	}
	fmt.Println("]")
}

func (parser *Parser) BuildSLRStates() (*Parser, error) {
    states := make([]*utils.LRState, 0)
    stateMap := make(map[string]*utils.LRState) // Usar un mapa para identificar estados únicos
	parser.SLR = &automatas.SLR{}

    // Inicializar el estado inicial con el cierre del primer ítem de la producción inicial
    initialStateItems := []*utils.Item{{
        Production: parser.Productions[0], // Asume que la primera producción es la producción inicial
        Position:   0,
        SubPos:     0,
    }}
    initialState := &utils.LRState{
        ID: 0,
        Items: parser.SLR.Closure(initialStateItems, parser.Productions),
        Transitions: make(map[string]int),
    }

	parser.SLR.StartState = initialState.ID

    states = append(states, initialState)
    stateMap[itemKeyForState(initialState.Items)] = initialState
    worklist := []*utils.LRState{initialState}

    nextID := 1

    for len(worklist) > 0 {
        currentState := worklist[0]
        worklist = worklist[1:]

        // Encuentra todos los símbolos después del punto en los ítems del estado actual
        symbols := findAllSymbols(currentState.Items)

        // Para cada símbolo, calcula el GOTO y verifica si el estado resultante ya existe
        for symbol := range symbols {
            newStateItems := parser.SLR.Goto(currentState.Items, symbol, parser.Productions)
            if len(newStateItems) == 0 {
                continue
            }

            key := itemKeyForState(newStateItems)
            newState, exists := stateMap[key]
            if !exists {
                newState = &utils.LRState{
                    ID: nextID,
                    Items: newStateItems,
                    Transitions: make(map[string]int),
                }
                nextID++
                states = append(states, newState)
                stateMap[key] = newState
                worklist = append(worklist, newState)
            }
            currentState.Transitions[symbol] = newState.ID
        }
    }
	parser.SLR.States = states

    return parser, nil
}

func findAllSymbols(items []*utils.Item) map[string]struct{} {
    symbols := make(map[string]struct{})
    for _, item := range items {
        if item.SubPos < len(item.Production.Body[item.Position]) {
            symbol := item.Production.Body[item.Position][item.SubPos]
            symbols[symbol] = struct{}{}
        }
    }
    return symbols
}

func itemKeyForState(items []*utils.Item) string {
    // Crear un slice de strings para mantener las claves de los ítems individuales
    itemKeys := make([]string, len(items))
    
    // Generar la clave para cada ítem
    for i, item := range items {
        // Asume que la posición del punto y la producción son suficientes para identificar el ítem
        // Puedes ajustar esta clave si es necesario incluir más información
        itemKeys[i] = fmt.Sprintf("%s-%d-%d", item.Production.Head, item.Position, item.SubPos)
    }
    
    // Ordenar las claves para garantizar la consistencia
    sort.Strings(itemKeys)
    
    // Concatenar todas las claves en una clave única para el estado
    return strings.Join(itemKeys, "|")
}


func (parser *Parser) First(symbol string, visited map[string]bool) map[string]struct{} {
    if visited[symbol] { // Chequeo de ciclo
        return make(map[string]struct{})
    }
    visited[symbol] = true
    defer delete(visited, symbol)

    set := make(map[string]struct{})
    if utils.IsTerminal(symbol) {
        set[symbol] = struct{}{}
        return set
    }

    for _, production := range parser.Productions {
        if production.Head == symbol {
            for _, body := range production.Body { // body es []string
                if len(body) > 0 {
                    firstSym := body[0] // Asume que quieres el primer símbolo de cada producción
                    if firstSym != symbol {
                        result := parser.First(firstSym, visited)
                        for s := range result {
                            set[s] = struct{}{}
                        }
                    }
                    // Suponiendo que la función debe detenerse al encontrar el primer símbolo no vacío
                    break
                }
            }
        }
    }
    return set
}



// Follow calcula el conjunto FOLLOW para un no terminal.
func (p *Parser) Follow(symbol string) []string {
    followSet := []string{}
    if symbol == p.Productions[0].Head {
        followSet = append(followSet, "$") // Añadir símbolo de fin de archivo si es la cabeza de la producción inicial
    }

    for _, prod := range p.Productions {
        for _, body := range prod.Body {
            for i, sym := range body {
                if sym == symbol {
                    if i+1 < len(body) {
                        // Llamar First para el siguiente símbolo en el cuerpo
                        firstOfNext := p.First(body[i+1], make(map[string]bool))
                        // Convertir el mapa a slice para usarlo en Filter y Contains
                        firstOfNextSlice := mapKeysToStringSlice(firstOfNext)
                        // Filtrar y añadir a followSet todos los símbolos excepto el épsilon
                        followSet = append(followSet, utils.Filter(firstOfNextSlice, func(s string) bool { return s != "ε" })...)
                        // Si épsilon está en firstOfNext, añadir el Follow del cabeza de la producción
                        if utils.Contains(firstOfNextSlice, "ε") {
                            followSet = append(followSet, p.Follow(prod.Head)...)
                        }
                    } else if i+1 == len(body) { // Si el símbolo es el último en el cuerpo, añadir el Follow del cabeza de la producción
                        followSet = append(followSet, p.Follow(prod.Head)...)
                    }
                }
            }
        }
    }
    // Eliminar duplicados en followSet
    return utils.Unique(followSet)
}

// Función auxiliar para convertir map[string]struct{} a []string
func mapKeysToStringSlice(m map[string]struct{}) []string {
    result := make([]string, 0, len(m))
    for key := range m {
        result = append(result, key)
    }
    return result
}

func (p *Parser) findProduction(symbol string) *utils.ProductionToken {
    for _, prod := range p.Productions {
        if prod.Head == symbol {
            return prod
        }
    }
    return nil
}


func (parser *Parser) ClosureLR1(items []*utils.Item, productions []*utils.ProductionToken) []*utils.Item {

	itemSet := make(map[string]*utils.Item)
    
    // Inicializar el conjunto con los ítems iniciales y sus lookaheads
    for _, item := range items {
        key := itemKeyLR1(item)
        itemSet[key] = item
    }

    changed := true
    for changed {
        changed = false
        currentItems := []*utils.Item{}
        for _, item := range itemSet {
            currentItems = append(currentItems, item)
        }

        for _, item := range currentItems {
            if item.SubPos < len(item.Production.Body[item.Position]) {
                nextSymbol := item.Production.Body[item.Position][item.SubPos]
                for _, prod := range productions {
                    if prod.Head == nextSymbol {
                        // Calcular los nuevos lookaheads para este nuevo ítem
                        newLookaheads := parser.calculateLookaheads(item, prod)
                        newItem := utils.NewLR1Item(prod, 0, 0, newLookaheads)
                        key := itemKeyLR1(newItem)
                        if _, exists := itemSet[key]; !exists {
                            itemSet[key] = newItem
                            changed = true
                        }
                    }
                }
            }
        }
    }

    // Convertir mapa a slice
    finalItems := make([]*utils.Item, 0, len(itemSet))
    for _, item := range itemSet {
        finalItems = append(finalItems, item)
    }
    return finalItems
}

func (parser *Parser) calculateLookaheads(item *utils.Item, prod *utils.ProductionToken) []string {
    result := []string{}
    remainingSymbols := item.Production.Body[item.Position][item.SubPos+1:]
    visited := make(map[string]bool)
    firstOfBeta := parser.CalculateFirstSequence(remainingSymbols, visited)

    for symbol := range firstOfBeta {
        if symbol != "ε" {
            result = append(result, symbol)
        }
    }

    if _, exists := firstOfBeta["ε"]; exists {
        result = append(result, item.Lookaheads...)
    }

    return utils.Unique(result)
}


func (parser *Parser) CalculateFirstSequence(sequence []string, visited map[string]bool) map[string]struct{} {
    result := make(map[string]struct{})
    if len(sequence) == 0 {
        // Si la secuencia es vacía, añadir ε
        result["ε"] = struct{}{}
        return result
    }

    // Calcula FIRST del primer símbolo de la secuencia
    firstOfFirst := parser.First(sequence[0], visited)

    // Añadir FIRST del primer símbolo al resultado
    for symbol := range firstOfFirst {
        result[symbol] = struct{}{}
        if symbol == "ε" && len(sequence) > 1 {
            // Si ε está en el FIRST del primer símbolo y hay más símbolos, calcular FIRST del resto de la secuencia
            restFirst := parser.CalculateFirstSequence(sequence[1:], visited)
            for sym := range restFirst {
                result[sym] = struct{}{}
            }
        }
    }

    return result
}

type LRTable struct {
    States      []*utils.LRState			`json:"states"`
    Transitions map[int]map[string]int  	`json:"transitions"`
    Actions     map[int]map[string]string 	`json:"actions"`
	Productions []*utils.ProductionToken 	`json:"productions"`
    Gotos       map[int]map[string]int   	`json:"gotos"`
}
func (table *LRTable) PrintTable(){
	fmt.Println("Estados:")
	for _,state := range table.States{
		fmt.Printf("Estado %d \n", state.ID)
	}
	fmt.Println("Transition:")
	for index,transition := range table.Transitions{
		for key, state := range transition{
			fmt.Printf("Transition %d with %s to state %d\n", index,key, state)
			
		}
	}
	fmt.Println("Acciones:")
	for index,accion := range table.Actions{
		for key, val := range accion{
			fmt.Printf("Action %d with %s action to %s\n", index,key, val)
			
		}
	}
	fmt.Println("Producciones:")
	for _,prod := range table.Productions{
		fmt.Printf("Production head: %s, body: %s\n", prod.Head, prod.Body)
	}
	fmt.Println("Gotos:")
	for index,accion := range table.Gotos{
		for key, val := range accion{
			fmt.Printf("Goto %d with %s action to %d state\n", index,key, val)
			
		}
	}
}

func (parser *Parser) allSymbols() []string {
    symbolSet := make(map[string]struct{})

    // Recolectar símbolos de las cabezas de producción
    for _, prod := range parser.Productions {
        symbolSet[prod.Head] = struct{}{}
        for _, body := range prod.Body {
            for _, sym := range body {
                symbolSet[sym] = struct{}{}
            }
        }
    }

    // Convertir el mapa a slice
    var symbols []string
    for sym := range symbolSet {
        symbols = append(symbols, sym)
    }

    sort.Strings(symbols) // opcional, para tener un orden consistente
    return symbols
}

// func (parser *Parser) BuildLR1Table() (*LRTable, error) {
//     states := []*utils.LRState{}
//     stateMap := make(map[string]int) // Para verificar existencia de estados por clave única.
//     table := &LRTable{
//         States:      states,
//         Transitions: make(map[int]map[string]int),
//         Actions:     make(map[int]map[string]string),
// 		Productions: parser.Productions,
//         Gotos:       make(map[int]map[string]int),
//     }

//     // Inicializar el estado inicial y agregarlo a la lista y mapa.
//     initialState := parser.ClosureLR1([]*utils.Item{{Production: parser.Productions[0], Position: 0, SubPos: 0, Lookaheads: []string{"$"}}}, parser.Productions)
//     initialStateKey := itemKeyForLR1State(initialState)
//     stateID := 0
//     states = append(states, &utils.LRState{ID: stateID, Items: initialState})
//     stateMap[initialStateKey] = stateID
//     worklist := []*utils.LRState{states[0]}

//     // Procesamiento de los estados
//     for len(worklist) > 0 {
//         currentState := worklist[0]
//         worklist = worklist[1:] // Dequeue

// 		if _, ok := table.Gotos[currentState.ID]; !ok {
//             table.Gotos[currentState.ID] = make(map[string]int)
//         }

//         currentTransitions := make(map[string]int)
//         currentActions := make(map[string]string)

//         allSymbols := parser.allSymbols() // Obtiene todos los símbolos
//         for _, symbol := range allSymbols {
//             newStateItems := parser.GotoLR1(currentState.Items, symbol)
//             if len(newStateItems) > 0 {
//                 newStateKey := itemKeyForLR1State(newStateItems)
//                 newStateID, exists := stateMap[newStateKey]
//                 if !exists {
//                     newStateID = len(states)
//                     states = append(states, &utils.LRState{ID: newStateID, Items: newStateItems})
//                     stateMap[newStateKey] = newStateID
//                     worklist = append(worklist, states[newStateID]) // Enqueue new state
//                 }
//                 // Agregar transiciones y acciones
//                 currentTransitions[symbol] = newStateID
//                 if utils.IsTerminal(symbol) {
//                     currentActions[symbol] = fmt.Sprintf("shift to %d", newStateID)
//                 } else {
//                     currentActions[symbol] = fmt.Sprintf("goto %d", newStateID)
// 					table.Gotos[currentState.ID][symbol] = newStateID
//                 }
//             }
//         }

//         table.Transitions[currentState.ID] = currentTransitions
//         table.Actions[currentState.ID] = currentActions
//     }

//     table.States = states // Asegurar que los estados están actualizados
//     return table, nil
// }

func (parser *Parser) BuildLR1Table() (*LRTable, error) {
    // Inicialización de la tabla y las estructuras de datos necesarias
    states := []*utils.LRState{}
    stateMap := make(map[string]int) // Mapa para controlar la existencia de estados
    table := &LRTable{
        States:      states,
        Transitions: make(map[int]map[string]int),
        Actions:     make(map[int]map[string]string),
        Productions: parser.Productions,
        Gotos:       make(map[int]map[string]int),
    }

    // Crear el estado inicial y procesar la clausura de la producción inicial
    initialState := parser.ClosureLR1([]*utils.Item{{Production: parser.Productions[0], Position: 0, SubPos: 0, Lookaheads: []string{"$"}}}, parser.Productions)
    initialStateKey := itemKeyForLR1State(initialState)
    stateID := 0
    states = append(states, &utils.LRState{ID: stateID, Items: initialState})
    stateMap[initialStateKey] = stateID
    worklist := []*utils.LRState{states[0]}

    // Bucle para procesar cada estado en la lista de trabajo
    for len(worklist) > 0 {
        currentState := worklist[0]
        worklist = worklist[1:] // Desencolar

        if _, ok := table.Gotos[currentState.ID]; !ok {
            table.Gotos[currentState.ID] = make(map[string]int)
        }

        currentTransitions := make(map[string]int)
        currentActions := make(map[string]string)

        // Obtener todos los símbolos a procesar
        allSymbols := parser.allSymbols()
        for _, symbol := range allSymbols {
            newStateItems := parser.GotoLR1(currentState.Items, symbol)
            if len(newStateItems) > 0 {
                newStateKey := itemKeyForLR1State(newStateItems)
                newStateID, exists := stateMap[newStateKey]
                if !exists {
                    newStateID = len(states)
                    states = append(states, &utils.LRState{ID: newStateID, Items: newStateItems})
                    stateMap[newStateKey] = newStateID
                    worklist = append(worklist, states[newStateID]) // Encolar nuevo estado
                }
                // Definir transiciones y acciones
                currentTransitions[symbol] = newStateID
                if utils.IsTerminal(symbol) {
                    currentActions[symbol] = fmt.Sprintf("shift to %d", newStateID)
                } else {
                    currentActions[symbol] = fmt.Sprintf("goto %d", newStateID)
                    table.Gotos[currentState.ID][symbol] = newStateID
                }
            }
        }

        table.Transitions[currentState.ID] = currentTransitions
        table.Actions[currentState.ID] = currentActions
    }

    table.States = states // Actualizar la lista de estados
    return table, nil
}





func itemKeyForLR1State(items []*utils.Item) string {
    // Genera una clave única para un estado basado en sus ítems
    // Debes ajustar esto para incluir los lookaheads en la clave
    var keys []string
    for _, item := range items {
        keys = append(keys, fmt.Sprintf("%s-%d-%d-%s", item.Production.Head, item.Position, item.SubPos, strings.Join(item.Lookaheads, ",")))
    }
    sort.Strings(keys)
    return strings.Join(keys, "|")
}


func ReverseAcceptedExpArray(expArray []*automatas.AcceptedExp) {
	for i, j := 0, len(expArray)-1; i < j; i, j = i+1, j-1 {
		expArray[i], expArray[j] = expArray[j], expArray[i]
	}
}


func (parser *Parser) ParseString(input string, table *LRTable) (bool, error) {
    rawTokens, err := parser.scanner.ScanFile(input)
    if err != nil {
        return false, err
    }

    // Filtrar tokens ignorados
    tokens := []*automatas.AcceptedExp{}
    for _, token := range rawTokens {
        if !parser.isIgnoredToken(token.Token) {
            tokens = append(tokens, token)
        }
    }

    ReverseAcceptedExpArray(tokens)

	table.PrintTable()

    tokens = append(tokens, &automatas.AcceptedExp{Token: "$"})  // Añadir el token de EOF

    stateStack := []int{0}  // La pila de estados empieza con el estado inicial
    symbolStack := []string{}  // Pila de símbolos

    index := 0
    for index < len(tokens) {
        currentState := stateStack[len(stateStack)-1]
        currentToken := tokens[index].Token

		fmt.Printf("Token %s with value %s\n",currentToken, tokens[index].Value)


        actions, exists := table.Actions[currentState]
        if !exists {
            return false, fmt.Errorf("no actions found for state: %d", currentState)
        }

        action, ok := actions[currentToken]
        if !ok {
            return false, fmt.Errorf("no action for token %s in state %d", currentToken, currentState)
        }

        if action[:1] == "s" {  // Shift action
            newStateStr := action[len(action)-1:]
            newState, err := strconv.Atoi(newStateStr)
			fmt.Println(newStateStr)
            if err != nil {
                return false, fmt.Errorf("invalid state number: %s", newStateStr)
            }
            stateStack = append(stateStack, newState)
            symbolStack = append(symbolStack, currentToken)
            index++  // Avanzar al siguiente token
        } else if action[:1] == "r" {  // Reduce action
            prodIndexStr := action[2:]
            prodIndex, err := strconv.Atoi(prodIndexStr)
            if err != nil {
                return false, fmt.Errorf("invalid production index: %s", prodIndexStr)
            }

            production := table.Productions[prodIndex]
            if len(stateStack) < len(production.Body) {
                return false, fmt.Errorf("stack has fewer elements than the production body")
            }

            // Pop the stack by the length of the production body
            stateStack = stateStack[:len(stateStack)-len(production.Body)]
            symbolStack = symbolStack[:len(symbolStack)-len(production.Body)]

            // Push the nonterminal onto the symbol stack
            symbolStack = append(symbolStack, production.Head)

            // Use the GOTO table to find the next state
            gotoState, exists := table.Gotos[stateStack[len(stateStack)-1]][production.Head]
            if !exists {
                return false, fmt.Errorf("no goto entry for nonterminal %s in state %d", production.Head, stateStack[len(stateStack)-1])
            }
            stateStack = append(stateStack, gotoState)
        } else if action == "accept" {
            return true, nil
        }
    }

    return false, fmt.Errorf("input did not resolve to an accept state")
}

// isIgnoredToken verifica si un token debe ser ignorado según la lista de IgnoredTokens del parser.
func (parser *Parser) isIgnoredToken(token string) bool {
    for _, ignoredToken := range parser.IgnoredTokens {
        if ignoredToken.Token == token {
            return true
        }
    }
    return false
}



func contains(tokens []utils.ParseToken, token string) bool {
    for _, t := range tokens {
        if t.Token == token {
            return true
        }
    }
    return false
}



func itemKeyLR1(item *utils.Item) string {
    // Clave única para ítem LR(1) que incluye los lookaheads
    return fmt.Sprintf("%s-%d-%d-%v", item.Production.Head, item.Position, item.SubPos, item.Lookaheads)
}

func (parser *Parser) GotoLR1(items []*utils.Item, symbol string) []*utils.Item {
    movedItems := []*utils.Item{}
    for _, item := range items {
        if item.SubPos < len(item.Production.Body[item.Position]) && item.Production.Body[item.Position][item.SubPos] == symbol {
            // Crear un nuevo ítem moviendo el punto pasado el símbolo
            newItem := &utils.Item{
                Production: item.Production,
                Position: item.Position,
                SubPos: item.SubPos + 1,
                Lookaheads: item.Lookaheads, // Los lookaheads permanecen
            }
            movedItems = append(movedItems, newItem)
        }
    }
    return parser.ClosureLR1(movedItems, parser.Productions) // Recalcular la clausura para los nuevos ítems
}


type SLRStateOutput struct {
	StateID   int       `json:"state_id"`
	Items     []SLRItem `json:"items"`
	Actions   []string  `json:"actions"`
}

type SLRItem struct {
	ProductionHead string   `json:"production_head"`
	ProductionBody []string `json:"production_body"`
	Description    string   `json:"description"`
}


func (parser *Parser) PrintSLR() ([]byte, error) {
	slr := parser.SLR
	output := []SLRStateOutput{}

	for _, state := range slr.States {
		stateOutput := SLRStateOutput{
			StateID: state.ID,
			Items:   []SLRItem{},
			Actions: []string{},
		}

		for _, item := range state.Items {
			if item.SubPos == len(item.Production.Body[item.Position]) { // Punto al final de la producción
				followSet := parser.Follow(item.Production.Head)
				description := fmt.Sprintf("[Reduce: %s → %v en %v]", item.Production.Head, item.Production.Body[item.Position], followSet)
				stateOutput.Actions = append(stateOutput.Actions, description)
			} else {
				nextSymbol := item.Production.Body[item.Position][item.SubPos]
				description := fmt.Sprintf("[Item: %s → %v]", item.Production.Head, item.Production.Body[item.Position])
				stateOutput.Items = append(stateOutput.Items, SLRItem{
					ProductionHead: item.Production.Head,
					ProductionBody: item.Production.Body[item.Position],
					Description:    description,
				})
				if state.Transitions != nil {
					nextState := state.Transitions[nextSymbol]
					action := fmt.Sprintf("[Desplazar: %s a Estado %d]", nextSymbol, nextState)
					stateOutput.Actions = append(stateOutput.Actions, action)
				}
			}
		}

		output = append(output, stateOutput)
	}

	// Serializar el resultado a JSON
	jsonData, err := json.Marshal(output)
	if err != nil {
		return nil, err
	}

	return jsonData, nil
}


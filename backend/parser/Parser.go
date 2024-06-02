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
				// fmt.Print("error parsing yapar: invalid token definition after splitter")
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
						return nil, fmt.Errorf("error parsing yapar: invalid token definition in token definitions")
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
			return nil, fmt.Errorf("error parsing yapar: invalid token definition twodots")

		case "SEMICOLON":
			// fmt.Println("semicolon found")
			return nil, fmt.Errorf("error parsing yapar: invalid token definition semicolon")

		case "OR":
			// fmt.Println("lowercase found")
			return nil, fmt.Errorf("error parsing yapar: invalid token definition or")

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
    // lex.scanner.PrintScanner()
	for _, item := range lex.scanner.Tokens {
        // fmt.Println(item.Token)
		if _, found := elementMap[item.Token]; !found {
			// Si un elemento de arr1 no está en arr2, retornar un error.
            // fmt.Println(item.Token)
			return nil, fmt.Errorf("token in yapar not in yalex")
		}
	}

	// fmt.Print(lex.Productions)

    // for _, prod := range lex.Productions{
    //     // fmt.Printf("head: %s, body: %s \n", prod.Head, prod.Body)
    // }

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
    initialStateItems := parser.Closure([]*utils.Item{{
        Production: parser.Productions[0], // Asume que la primera producción es la producción inicial
        Position:   0,
        SubPos:     0,
    }}, parser.Productions)

    initialState := &utils.LRState{
        ID: 0,
        Items: initialStateItems,
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
            newStateItems := parser.Closure(parser.SLR.Goto(currentState.Items, symbol, parser.Productions), parser.Productions)
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

func (parser *Parser) Closure(items []*utils.Item, productions []*utils.ProductionToken) []*utils.Item {
    closure := make([]*utils.Item, len(items))
    copy(closure, items) // Copia los ítems iniciales al cierre

    added := true
    for added {
        added = false
        newItems := []*utils.Item{}

        for _, item := range closure {
            if item.SubPos < len(item.Production.Body[item.Position]) {
                symbol := item.Production.Body[item.Position][item.SubPos]
                // Verificar si el símbolo es un no terminal
                if !utils.IsTerminal(symbol) {
                    // Agregar todas las producciones que empiezan con este no terminal
                    for _, prod := range productions {
                        if prod.Head == symbol {
                            // Agregar cada producción posible para el no terminal
                            for _, body := range prod.Body {
                                newItem := &utils.Item{
                                    Production: &utils.ProductionToken{
                                        Head: prod.Head,
                                        Body: [][]string{body},
                                    },
                                    Position: 0,
                                    SubPos: 0,
                                }
                                if !containsItem(closure, newItem) && !containsItem(newItems, newItem) {
                                    newItems = append(newItems, newItem)
                                    added = true
                                }
                            }
                        }
                    }
                }
            }
        }

        // Agregar nuevos ítems al cierre
        closure = append(closure, newItems...)
    }

    return closure
}

// Helper function to check if an item already exists in a slice of items
func containsItem(items []*utils.Item, item *utils.Item) bool {
    for _, itm := range items {
        if itm.Production.Head == item.Production.Head && len(itm.Production.Body) == len(item.Production.Body) {
            match := true
            for i := range itm.Production.Body {
                if !equalBodies(itm.Production.Body[i], item.Production.Body[i]) {
                    match = false
                    break
                }
            }
            if match {
                return true
            }
        }
    }
    return false
}

func equalBodies(a, b []string) bool {
    if len(a) != len(b) {
        return false
    }
    for i := range a {
        if a[i] != b[i] {
            return false
        }
    }
    return true
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




func Unique(items []string) []string {
    seen := make(map[string]struct{})
    result := []string{}
    for _, item := range items {
        if _, ok := seen[item]; !ok {
            seen[item] = struct{}{}
            result = append(result, item)
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
func areSlicesEqual(a, b []string) bool {
    if len(a) != len(b) {
        return false
    }
    for i := range a {
        if a[i] != b[i] {
            return false
        }
    }
    return true
}


func (parser *Parser) ClosureLR1(items []*utils.Item, productions []*utils.ProductionToken) []*utils.Item {
    closure := append([]*utils.Item{}, items...)

    added := true
    for added {
        added = false
        newItems := []*utils.Item{}

        for _, item := range closure {
            if item.SubPos < len(item.Production.Body[item.Position]) {
                nextSymbol := item.Production.Body[item.Position][item.SubPos]
                if !utils.IsTerminal(nextSymbol) {
                    nextSymbolLookaheads := []string{}
                    if item.SubPos+1 < len(item.Production.Body[item.Position]) {
                        nextNextSymbol := item.Production.Body[item.Position][item.SubPos+1]
                        lookaheadResults := parser.First(nextNextSymbol, make(map[string]bool))
                        nextSymbolLookaheads = setToStringSlice(lookaheadResults)
                    } else {
                        nextSymbolLookaheads = item.Lookaheads
                    }

                    for _, prod := range productions {
                        if prod.Head == nextSymbol {
                            for _, body := range prod.Body {
                                newItem := &utils.Item{
                                    Production: &utils.ProductionToken{
                                        Head: prod.Head,
                                        Body: [][]string{body},
                                    },
                                    Position: 0,
                                    SubPos: 0,
                                    Lookaheads: nextSymbolLookaheads,
                                }
                                if !containsItemLR1(closure, newItem) {
                                    newItems = append(newItems, newItem)
                                    added = true
                                }
                            }
                        }
                    }
                }
            }
        }
        closure = append(closure, newItems...)
    }
    return closure
}



func (parser *Parser) GotoLR1(items []*utils.Item, symbol string, productions []*utils.ProductionToken) []*utils.Item {
    movedItems := []*utils.Item{}
    for _, item := range items {
        if item.SubPos < len(item.Production.Body[item.Position]) && item.Production.Body[item.Position][item.SubPos] == symbol {
            newItem := &utils.Item{
                Production: item.Production,
                Position: item.Position,
                SubPos: item.SubPos + 1,
                Lookaheads: item.Lookaheads,
            }
            movedItems = append(movedItems, newItem)
        }
    }
    return parser.ClosureLR1(movedItems, productions)
}

func setToStringSlice(set map[string]struct{}) []string {
    var slice []string
    for key := range set {
        slice = append(slice, key)
    }
    return slice
}


func containsItemLR1(items []*utils.Item, item *utils.Item) bool {
    for _, itm := range items {
        if itm.Production.Head == item.Production.Head && len(itm.Production.Body) == len(item.Production.Body) {
            match := true
            for i := range itm.Production.Body {
                if itm.Production.Head == item.Production.Head && areSlicesEqual(itm.Production.Body[i], item.Production.Body[i]) && equalLookaheads(itm.Lookaheads, item.Lookaheads) {
                    match = true
                    break
                }                
            }
            if match {
                return true
            }
        }
    }
    return false
}

func equalLookaheads(a, b []string) bool {
    if len(a) != len(b) {
        return false
    }
    sort.Strings(a)
    sort.Strings(b)
    for i := range a {
        if a[i] != b[i] {
            return false
        }
    }
    return true
}

func (parser *Parser) BuildLR1States() ([]*utils.LRState, error) {
    initialState := parser.ClosureLR1([]*utils.Item{
        { // Asume que la primera producción es la producción extendida S' -> S
            Production: parser.Productions[0],
            Position:   0,
            SubPos:     0,
            Lookaheads: []string{"$"}, // EOF symbol
        },
    }, parser.Productions)

    states := []*utils.LRState{{Items: initialState}}
    stateMap := map[string]int{itemKeyForLR1State(initialState): 0}
    worklist := []int{0}

    nextID := 1

    for len(worklist) > 0 {
        stateID := worklist[0]
        worklist = worklist[1:]
        currentState := states[stateID]

        // Find all symbols after the dot in any item of the state
        symbols := findAllSymbols(currentState.Items)
        for symbol := range symbols {
            newStateItems := parser.GotoLR1(currentState.Items, symbol, parser.Productions)
            if len(newStateItems) == 0 {
                continue
            }

            stateKey := itemKeyForLR1State(newStateItems)
            newStateID, exists := stateMap[stateKey]
            if !exists {
                newStateID = nextID
                nextID++
                states = append(states, &utils.LRState{ID: newStateID, Items: newStateItems})
                stateMap[stateKey] = newStateID
                worklist = append(worklist, newStateID)
            }
            // Update transition map
            if currentState.Transitions == nil {
                currentState.Transitions = make(map[string]int)
            }
            currentState.Transitions[symbol] = newStateID
        }
    }

    return states, nil
}

// func (parser *Parser) BuildLR1Table(states []*utils.LRState) (*LRTable, error) {
//     table := &LRTable{
//         States:      states,
//         Transitions: make(map[int]map[string]int),
//         Actions:     make(map[int]map[string]string),
//         Gotos:       make(map[int]map[string]int),
//         Productions: parser.Productions,
//     }

//     for _, state := range states {
//         table.Actions[state.ID] = make(map[string]string)
//         table.Gotos[state.ID] = make(map[string]int)

//         for _, item := range state.Items {
//             nextSymbolIndex := item.SubPos
//             if nextSymbolIndex < len(item.Production.Body[item.Position]) {
//                 nextSymbol := item.Production.Body[item.Position][nextSymbolIndex]
//                 if utils.IsTerminal(nextSymbol) {
//                     nextState, exists := state.Transitions[nextSymbol]
//                     if exists {
//                         table.Actions[state.ID][nextSymbol] = fmt.Sprintf("s%d", nextState)
//                     }
//                 } else { // Non-terminal
//                     nextState, exists := state.Transitions[nextSymbol]
//                     if exists {
//                         table.Gotos[state.ID][nextSymbol] = nextState
//                     }
//                 }
//             } else if nextSymbolIndex == len(item.Production.Body[item.Position]) { // Reduce
//                 for _, lookahead := range item.Lookaheads {
//                     // Find the production index
//                     prodIndex := -1
//                     for i, prod := range parser.Productions {
//                         // fmt.Println(item.Production.Body, item.Production.Head)
//                         if prod.Head == item.Production.Head {
//                             prodIndex = i
//                             break
//                         }
//                     }
//                     if prodIndex == -1 {
//                         return nil, fmt.Errorf("production not found")
//                     }
//                     action := fmt.Sprintf("r%d", prodIndex)
//                     if item.Production.Head == parser.Productions[0].Head { // Accept
//                         action = "accept"
//                     }
//                     table.Actions[state.ID][lookahead] = action
//                 }
//             }
//         }
//     }

//     return table, nil
// }

func (parser *Parser) BuildLR1Table(states []*utils.LRState) (*LRTable, error) {
    table := &LRTable{
        States:      states,
        Transitions: make(map[int]map[string]int),
        Actions:     make(map[int]map[string]string),
        Gotos:       make(map[int]map[string]int),
        Productions: parser.Productions,
    }

    for _, state := range states {
        table.Actions[state.ID] = make(map[string]string)
        table.Gotos[state.ID] = make(map[string]int)

        for _, item := range state.Items {
            nextSymbolIndex := item.SubPos
            if nextSymbolIndex < len(item.Production.Body[item.Position]) {
                nextSymbol := item.Production.Body[item.Position][nextSymbolIndex]
                if utils.IsTerminal(nextSymbol) {
                    nextState, exists := state.Transitions[nextSymbol]
                    if exists {
                        table.Actions[state.ID][nextSymbol] = fmt.Sprintf("s%d", nextState)
                    }
                } else {
                    nextState, exists := state.Transitions[nextSymbol]
                    if exists {
                        table.Gotos[state.ID][nextSymbol] = nextState
                    }
                }
            } else if nextSymbolIndex == len(item.Production.Body[item.Position]) {
                for _, lookahead := range item.Lookaheads {
                    action := determineProductionAction(parser, item)
                    if action == "" {
                        return nil, fmt.Errorf("production not found")
                    }
                    if item.Production.Head == parser.Productions[0].Head && nextSymbolIndex == len(item.Production.Body[item.Position]) {
                        action = "accept"
                    }
                    table.Actions[state.ID][lookahead] = action
                }
            }
        }
    }

    return table, nil
}

func determineProductionAction(parser *Parser, item *utils.Item) string {
    for i, prod := range parser.Productions {
        if prod.Head == item.Production.Head {
            bodyIndex, found := findBodyIndex(prod.Body, item.Production.Body[item.Position])
            if found {
                return fmt.Sprintf("r|%d|%d|%d", i, bodyIndex, item.SubPos-1)
            }
        }
    }
    return ""
}

// Verifica si un conjunto de subproducciones contiene una subproducción específica
func findBodyIndex(bodies [][]string, targetBody []string) (int, bool) {
    for i, body := range bodies {
        if areSlicesEqual(body, targetBody) {
            return i, true
        }
    }
    return -1, false
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
    // symbolStack := []string{}  // Pila de símbolos

    index := 0
    for index < len(tokens) {
        currentState := stateStack[len(stateStack)-1]
        currentToken := tokens[index].Token

        actions, exists := table.Actions[currentState]
        if !exists {
            return false, fmt.Errorf("no actions found for state: %d", currentState)
        }

        action, ok := actions[currentToken]
        if !ok {
            return false, fmt.Errorf("no action for token %s in state %d", currentToken, currentState)
        }

        fmt.Printf("Current state: %d, Token: %s, Action: %s\n", currentState, currentToken, action)

        actionParts := strings.Split(action, "|")
        if actionParts[0][:1] == "s" {  // Shift action
            newStateStr := actionParts[0][1:]
            newState, err := strconv.Atoi(newStateStr)
            if err != nil {
                return false, fmt.Errorf("invalid state number: %s", newStateStr)
            }
            stateStack = append(stateStack, newState)
            // symbolStack = append(symbolStack, currentToken)
            index++  // Avanzar al siguiente token
        } else if actionParts[0][:1] == "r" {  // Reduce action
            prodIndex, _ := strconv.Atoi(actionParts[1])
            bodyIndex, _ := strconv.Atoi(actionParts[2])
            subPosIndex, _ := strconv.Atoi(actionParts[3])
            production := table.Productions[prodIndex]
            
            fmt.Printf("StateStack: %v, Production: %v, Index: %d, Subpos: %d\n", stateStack, production.Body, bodyIndex, subPosIndex)

            // if len(symbolStack) < len(production.Body[bodyIndex]) {
            //     return false, fmt.Errorf("stack has fewer elements than the production body requires (%d needed, %d present)", len(production.Body[bodyIndex]), len(symbolStack))
            // }
            
            stateStack = stateStack[:len(stateStack)-len(production.Body[bodyIndex])]
            // symbolStack = symbolStack[:len(symbolStack)-len(production.Body[bodyIndex])]
            // symbolStack = append(symbolStack, production.Head)
        
            if len(stateStack) == 0 {
                return false, fmt.Errorf("state stack is empty after reduction")
            }
            currentState := stateStack[len(stateStack)-1]
            gotoState, exists := table.Gotos[currentState][production.Head]
            if !exists {
                return false, fmt.Errorf("no goto entry for nonterminal %s in state %d", production.Head, currentState)
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


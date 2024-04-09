package automatas

import (
	"encoding/json"
	"fmt"
	"github.com/adrianfulla/compiler/backend/utils"
)

func JsonToDfa(jsonData []byte) (*DAfdJson, error) {
	var dAfd DAfdJson
	err := json.Unmarshal(jsonData, &dAfd)
	if err != nil {
		return nil, err
	}
	return &dAfd, nil
}

func SimulateDFA(expresion string, estadoInicial string, estadosFinales []string, transiciones map[string]map[string]string) (int, string, error) {
	fmt.Printf("\nSimulate DFA\n")
	currState := estadoInicial
	var lastAcceptedState string
	var lastAcceptedCount int
	if expresion != "" {
		for count, char := range expresion {
			fmt.Printf("Simulate %d, %s, %s\n", count, string(char), currState)
			next_state, err := moveState(currState, char, transiciones)
			if err != nil {
				fmt.Printf("Returning error, transition not found\n")
				return count, currState, err
			}
			currState = next_state
			if utils.StringInStringArray(currState, estadosFinales) {
				if count != len(expresion)-1 {
					next_char := []rune(expresion[count+1 : count+2])
					fmt.Print(string(next_char))
					next_state, err := moveState(currState, next_char[0], transiciones)
					if err != nil {
						fmt.Printf("Returning true\n")
						return count, currState, nil
					}
					lastAcceptedCount = count
					lastAcceptedState = currState
					count++
					currState = next_state
				} else {
					fmt.Printf("Returning true\n")
					return count, currState, nil
				}
			}
		}
	}
	fmt.Printf("Returning partial true\n")
	return lastAcceptedCount, lastAcceptedState, nil
}

func moveState(curr_state string, char rune, transiciones map[string]map[string]string) (string, error) {
	for state, transicion := range transiciones {
		if state == curr_state {
			for sim, next_state := range transicion {
				if sim == string(char) {
					fmt.Printf("Move %s, %s, %s\n", state, sim, next_state)
					return next_state, nil
				}
			}
		}
	}
	return curr_state, fmt.Errorf("error in moveState: no transition for %s with %s symbol", curr_state, string(char))
}

func ExtendedSimulateAfd(expresion string, afd DAfdJson) {
	acceptedStack := utils.Stack{}
	traveled := 0

	for traveled < len(expresion)-1 {
		fmt.Printf("\nTraveled:%d, %d", traveled, len(expresion))
		tryExp := expresion[traveled:]
		returnCount, curr_state, err := SimulateDFA(tryExp, afd.EstadoInicial, afd.EstadosFinales, afd.Transiciones)

		if err == nil {
			fmt.Printf("Expresion valida encontrada %s desde %d hasta %d con %s\n", "", traveled, traveled+returnCount, curr_state)
			acceptedStack.Push(&AcceptedExp{
				Start: traveled,
				End:   returnCount,
				Value: expresion[traveled : traveled+returnCount+1],
			})
		}
		traveled = traveled + returnCount + 1
	}

	fmt.Print("\nCadenas acceptadas\n")
	for acceptedStack.Size() > 0 {
		acceptedExp := acceptedStack.Pop().(*AcceptedExp)
		fmt.Printf("[Start: %d, End: %d, Value: %s]\n", acceptedExp.Start, acceptedExp.End, acceptedExp.Value)
	}
}

type AcceptedExp struct {
	Start int    `json:"start"`
	End   int    `json:"end"`
	Value string `json:"value"`
}

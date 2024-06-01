package automatas

import (
	"fmt"
	"github.com/adrianfulla/compiler/backend/utils"
)


type SLR struct{
	States []*utils.LRState `json:"states"`
	StartState int	`json:"start_state"`
}

func (slr SLR) Closure(items []*utils.Item, productions []*utils.ProductionToken) []*utils.Item {
    itemSet := make(map[string]*utils.Item)
    
    // Inicializar el conjunto con los ítems iniciales
    for _, item := range items {
        key := itemKey(item)
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
                nextSymbol := item.Production.Body[item.Position][item.SubPos] // Accede al símbolo específico
                for _, prod := range productions {
                    if prod.Head == nextSymbol {
                        newItem := &utils.Item{Production: prod, Position: 0, SubPos: 0}
                        key := itemKey(newItem)
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

// Utiliza una función de clave única para identificar ítems
func itemKey(item *utils.Item) string {
    return fmt.Sprintf("%s-%d-%d", item.Production.Head, item.Position, item.SubPos)
}


func contains(items []*utils.Item, newItem *utils.Item) bool {
    for _, item := range items {
        // Suponiendo que Production tiene un campo `ID` o puedes comparar los `Head`.
        if item.Production.Head == newItem.Production.Head && item.Position == newItem.Position && item.SubPos == newItem.SubPos {
            return true
        }
    }
    return false
}


func (slr SLR) Goto(items []*utils.Item, symbol string, productions []*utils.ProductionToken) []*utils.Item {
    movedItems := []*utils.Item{}

    // Mueve el punto si el símbolo coincide
    for _, item := range items {
        if item.SubPos < len(item.Production.Body[item.Position]) && item.Production.Body[item.Position][item.SubPos] == symbol {
            newItem := &utils.Item{Production: item.Production, Position: item.Position, SubPos: item.SubPos + 1}
            movedItems = append(movedItems, newItem)
        }
    }

    // Calcula el cierre del conjunto de ítems movidos
    return slr.Closure(movedItems, productions)
}



func (slr SLR) PrintSLR(){
	for _,state := range slr.States{
		fmt.Printf("Estado %d\n", state.ID)
		for _, item := range state.Items{
			fmt.Printf("Item %s con pos %d y subpos %d\n", item.Production.Head, item.Position, item.SubPos)
		}
		for key, transition := range state.Transitions{
			fmt.Printf("Transition con %s hacia estado %d\n",key, transition)
		}
	}
}

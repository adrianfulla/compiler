package utils

import (
	// "fmt"
	// "reflect"
	// "unicode"
)


type Item struct {
    Production *ProductionToken `json:"production"`
    Position   int              `json:"position"`
    SubPos     int              `json:"subposition"`
    Lookaheads []string         `json:"lookaheads"`
}

func NewLR1Item(prod *ProductionToken, pos int, subPos int, lookahead []string) *Item {
    return &Item{
        Production: prod,
        Position:   pos,
        SubPos:     subPos,
        Lookaheads: lookahead,
    }
}

type LRState struct {
    ID      int `json:"id"`
    Items []*Item   `json:"items"`
    Transitions map[string]int `json:"transitions"`
}

type ValidSegment struct {
    Start   int `json:"start"`
    End   int `json:"end"`
    Message   string `json:"message"`
}
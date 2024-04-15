package utils

import "fmt"

type Token struct {
	Name       string `json:"name"`
	Expression string `json:"expression"`
	Action     string `json:"Action"`
}

type RegexToken struct{
	Value []rune
	IsOperator string
}
func (r *RegexToken) String() string{
	return fmt.Sprintf("{Value: %s, ValueString:%s, IsOperator: %s}", r.Value, r.IsOperator)
}

type LexToken struct{
	Token   string `json:"token"`
	Regex 	string `json:"regex"`
	Action  string `json:"action"`
}


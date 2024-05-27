module github.com/adrianfulla/compiler/backend/parser

go 1.22.1

replace github.com/adrianfulla/compiler/backend/automatas => ../automatas

replace github.com/adrianfulla/compiler/backend/utils => ../utils

require (
	github.com/adrianfulla/compiler/backend/automatas v0.0.0-00010101000000-000000000000
	github.com/adrianfulla/compiler/backend/lexer v0.0.0-00010101000000-000000000000
	github.com/adrianfulla/compiler/backend/utils v0.0.0-00010101000000-000000000000
)

replace github.com/adrianfulla/compiler/backend/lexer => ../lexer

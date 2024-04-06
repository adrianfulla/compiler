module github.com/adrianfulla/compiler/backend

go 1.22.1

replace github.com/adrianfulla/compiler/backend/automatas => ./automatas

require github.com/adrianfulla/compiler/backend/automatas v0.0.0-00010101000000-000000000000

require github.com/adrianfulla/compiler/backend/utils v0.0.0-00010101000000-000000000000 // indirect

replace github.com/adrianfulla/compiler/backend/utils => ./utils

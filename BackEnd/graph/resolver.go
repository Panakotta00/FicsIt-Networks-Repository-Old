package graph

//go:generate go run github.com/99designs/gqlgen generate

import "FINRepository/graph/model"

type Resolver struct {
	todos []*model.Todo
}
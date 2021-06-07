package query

import "context"

type Query interface {
	Name() string
}

type QueryHandler interface {
	Handle(context.Context, Query) (interface{}, error)
}

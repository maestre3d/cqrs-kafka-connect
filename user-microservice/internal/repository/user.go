package repository

import (
	"context"

	"github.com/maestre3d/cqrs-kafka-connect/user-microservice/internal/aggregate"
)

type User interface {
	Transaction
	Save(context.Context, aggregate.User) error
	Update(context.Context, aggregate.User) error
	Find(context.Context, string) (*aggregate.User, error)
	Search(context.Context, string) ([]*aggregate.User, error)
}

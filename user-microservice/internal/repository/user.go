package repository

import (
	"context"

	"github.com/maestre3d/cqrs-kafka-connect/user-microservice/internal/aggregate"
)

type User interface {
	Save(context.Context, aggregate.User) error
	Find(context.Context, string) (*aggregate.User, error)
	Search(context.Context, string) ([]*aggregate.User, error)
}

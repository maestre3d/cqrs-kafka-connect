package repository

import (
	"context"

	"github.com/maestre3d/cqrs-kafka-connect/user-microservice/internal/aggregate"
	"github.com/maestre3d/cqrs-kafka-connect/user-microservice/internal/domain"
)

type User interface {
	domain.Transaction
	Save(context.Context, aggregate.User) error
	Update(context.Context, aggregate.User) error
	Find(context.Context, string) (*aggregate.User, error)
	Search(context.Context, domain.Criteria) ([]*aggregate.User, string, error)
}

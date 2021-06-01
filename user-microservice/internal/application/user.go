package application

import (
	"context"

	"github.com/maestre3d/cqrs-kafka-connect/user-microservice/internal/aggregate"
	"github.com/maestre3d/cqrs-kafka-connect/user-microservice/internal/repository"
)

type User struct {
	repo repository.User
}

func NewUser(r repository.User) *User {
	return &User{
		repo: r,
	}
}

func (u *User) Create(ctx context.Context, id, username, displayName string) error {
	user := aggregate.NewUser(id, username, displayName)
	return u.repo.Save(ctx, *user)
}

func (u *User) Search(ctx context.Context, query string) ([]*aggregate.User, error) {
	return u.repo.Search(ctx, query)
}

func (u *User) GetById(ctx context.Context, userID string) (*aggregate.User, error) {
	return u.repo.Find(ctx, userID)
}

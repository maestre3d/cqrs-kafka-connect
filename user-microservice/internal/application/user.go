package application

import (
	"context"
	"errors"

	"github.com/maestre3d/cqrs-kafka-connect/user-microservice/internal/aggregate"
	"github.com/maestre3d/cqrs-kafka-connect/user-microservice/internal/repository"
)

type User interface {
	UserReadOnly
	UserWriteOnly
}

type UserWriteOnly interface {
	Create(ctx context.Context, id, username, displayName string) error
	Update(ctx context.Context, id, displayName string) error
}

type UserReadOnly interface {
	Search(ctx context.Context, query string) ([]*aggregate.User, error)
	GetById(ctx context.Context, userID string) (*aggregate.User, error)
}

var (
	ErrUserNotFound = errors.New("user not found")
)

type user struct {
	repo repository.User
}

var _ User = &user{}

func (u *user) Create(ctx context.Context, id, username, displayName string) error {
	user := aggregate.NewUser(id, username, displayName)
	return u.repo.Save(ctx, *user)
}

func (u *user) Update(ctx context.Context, id, displayName string) error {
	user, err := u.repo.Find(ctx, id)
	if err != nil {
		return err
	} else if user == nil {
		return ErrUserNotFound
	}
	user.DisplayName = displayName
	return u.repo.Update(ctx, *user)
}

func (u *user) Search(ctx context.Context, query string) ([]*aggregate.User, error) {
	return u.repo.Search(ctx, query)
}

func (u *user) GetById(ctx context.Context, userID string) (*aggregate.User, error) {
	return u.repo.Find(ctx, userID)
}

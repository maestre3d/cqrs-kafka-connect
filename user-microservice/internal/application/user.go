package application

import (
	"context"
	"errors"

	"github.com/maestre3d/cqrs-kafka-connect/user-microservice/internal/aggregate"
	"github.com/maestre3d/cqrs-kafka-connect/user-microservice/internal/repository"
)

type User interface {
	Create(ctx context.Context, id, username, displayName string) error
	Update(ctx context.Context, id, displayName string) error
	Search(ctx context.Context, query string) ([]*aggregate.User, error)
	GetById(ctx context.Context, userID string) (*aggregate.User, error)
}

var (
	ErrUserNotFound = errors.New("user not found")
)

type UserImpl struct {
	repo repository.User
}

var _ User = &UserImpl{}

func (u *UserImpl) Create(ctx context.Context, id, username, displayName string) error {
	user := aggregate.NewUser(id, username, displayName)
	if err := u.repo.Save(ctx, *user); err != nil {
		return err
	}

	user.DisplayName = "primateArevalo"
	return u.repo.Update(ctx, *user)
}

func (u *UserImpl) Update(ctx context.Context, id, displayName string) error {
	user, err := u.repo.Find(ctx, id)
	if err != nil {
		return err
	} else if user == nil {
		return ErrUserNotFound
	}
	user.DisplayName = displayName
	return u.repo.Update(ctx, *user)
}

func (u *UserImpl) Search(ctx context.Context, query string) ([]*aggregate.User, error) {
	return u.repo.Search(ctx, query)
}

func (u *UserImpl) GetById(ctx context.Context, userID string) (*aggregate.User, error) {
	return u.repo.Find(ctx, userID)
}

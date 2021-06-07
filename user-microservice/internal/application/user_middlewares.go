package application

import (
	"context"

	"github.com/maestre3d/cqrs-kafka-connect/user-microservice/internal/aggregate"
	"github.com/maestre3d/cqrs-kafka-connect/user-microservice/internal/repository"
)

func NewUser(r repository.User) User {
	return &userRepositoryTransaction{
		repository: r,
		next: &user{
			repo: r,
		},
	}
}

func NewUserReadOnly(r repository.User) UserReadOnly {
	return &user{
		repo: r,
	}
}

type userRepositoryTransaction struct {
	repository repository.User
	next       User
}

var _ User = &userRepositoryTransaction{}

func (u *userRepositoryTransaction) Create(ctx context.Context, id, username, displayName string) (err error) {
	if err := u.repository.Begin(ctx); err != nil {
		return err
	}

	err = u.next.Create(ctx, id, username, displayName)
	commitOrRollback(u.repository, err)
	return
}

func (u *userRepositoryTransaction) Update(ctx context.Context, id, displayName string) (err error) {
	if err := u.repository.Begin(ctx); err != nil {
		return err
	}

	err = u.next.Update(ctx, id, displayName)
	commitOrRollback(u.repository, err)
	return
}

func (u *userRepositoryTransaction) Search(ctx context.Context, query string) ([]*aggregate.User, error) {
	return u.next.Search(ctx, query)
}

func (u *userRepositoryTransaction) GetById(ctx context.Context, userID string) (*aggregate.User, error) {
	return u.next.GetById(ctx, userID)
}

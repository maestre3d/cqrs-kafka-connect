package application

import (
	"context"

	"github.com/maestre3d/cqrs-kafka-connect/user-microservice/internal/aggregate"
	"github.com/maestre3d/cqrs-kafka-connect/user-microservice/internal/repository"
)

func NewUser(r repository.User) User {
	return &UserTransaction{
		repository: r,
		next: &UserImpl{
			repo: r,
		},
	}
}

type UserTransaction struct {
	repository repository.User
	next       User
}

var _ User = &UserTransaction{}

func (u *UserTransaction) commitOrRollback(err error) error {
	if err != nil {
		return u.repository.Rollback()
	}
	return u.repository.Commit()
}

func (u *UserTransaction) Create(ctx context.Context, id, username, displayName string) (err error) {
	if err := u.repository.Begin(ctx); err != nil {
		return err
	}

	err = u.next.Create(ctx, id, username, displayName)
	u.commitOrRollback(err)
	return
}

func (u *UserTransaction) Update(ctx context.Context, id, displayName string) (err error) {
	if err := u.repository.Begin(ctx); err != nil {
		return err
	}

	err = u.next.Update(ctx, id, displayName)
	u.commitOrRollback(err)
	return
}

func (u *UserTransaction) Search(ctx context.Context, query string) ([]*aggregate.User, error) {
	return u.next.Search(ctx, query)
}

func (u *UserTransaction) GetById(ctx context.Context, userID string) (*aggregate.User, error) {
	return u.next.GetById(ctx, userID)
}

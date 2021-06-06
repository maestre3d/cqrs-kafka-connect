package repository

import "context"

type Transaction interface {
	Begin(context.Context) error
	Commit() error
	Rollback() error
}

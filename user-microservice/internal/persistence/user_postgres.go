package persistence

import (
	"context"
	"database/sql"
	"sync"

	"github.com/maestre3d/cqrs-kafka-connect/user-microservice/internal/aggregate"
	"github.com/maestre3d/cqrs-kafka-connect/user-microservice/internal/repository"
)

type UserPostgres struct {
	db *sql.DB
	mu sync.RWMutex
}

func NewUserPostgres(db *sql.DB) *UserPostgres {
	return &UserPostgres{
		db: db,
		mu: sync.RWMutex{},
	}
}

var _ repository.User = &UserPostgres{}

func (u *UserPostgres) Save(ctx context.Context, user aggregate.User) error {
	u.mu.Lock()
	defer u.mu.Unlock()

	conn, err := u.db.Conn(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	tx, err := conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := `INSERT INTO users(user_id, username, display_name) VALUES ($1, $2, $3)`
	if _, err = tx.ExecContext(ctx, q, user.ID, user.Username, user.DisplayName); err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (u *UserPostgres) Update(ctx context.Context, user aggregate.User) error {
	u.mu.Lock()
	defer u.mu.Unlock()

	conn, err := u.db.Conn(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	tx, err := conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := `UPDATE users SET username = $1, display_name = $2 WHERE user_id = $3`
	if _, err = tx.ExecContext(ctx, q, user.ID, user.Username, user.DisplayName); err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (u *UserPostgres) Find(ctx context.Context, userID string) (*aggregate.User, error) {
	u.mu.RLock()
	defer u.mu.RUnlock()

	conn, err := u.db.Conn(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	q := `SELECT user_id, username, display_name FROM users WHERE user_id = $1`
	user := new(aggregate.User)
	err = conn.QueryRowContext(ctx, q, userID).Scan(&user.ID, &user.Username, &user.DisplayName)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserPostgres) Search(ctx context.Context, query string) ([]*aggregate.User, error) {
	return nil, nil // this will be implemented using full-text search through ElasticSearch
}

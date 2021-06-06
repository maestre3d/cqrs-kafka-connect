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
	tx *sql.Tx
	mu sync.RWMutex
}

func NewUserPostgres(db *sql.DB) *UserPostgres {
	return &UserPostgres{
		db: db,
		mu: sync.RWMutex{},
	}
}

var _ repository.User = &UserPostgres{}

func (u *UserPostgres) Begin(ctx context.Context) error {
	u.mu.Lock()
	defer u.mu.Unlock()

	var err error
	u.tx, err = u.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	return nil
}

func (u *UserPostgres) Commit() error {
	return u.tx.Commit()
}

func (u *UserPostgres) Rollback() error {
	return u.tx.Rollback()
}

func (u *UserPostgres) Save(ctx context.Context, user aggregate.User) error {
	u.mu.Lock()
	defer u.mu.Unlock()

	q := `INSERT INTO users(user_id, username, display_name) VALUES ($1, $2, $3)`
	stmt, err := u.tx.PrepareContext(ctx, q)
	if err != nil {
		return err
	}

	if _, err := stmt.ExecContext(ctx, user.ID, user.Username, user.DisplayName); err != nil {
		return err
	}
	return nil
}

func (u *UserPostgres) Update(ctx context.Context, user aggregate.User) error {
	u.mu.Lock()
	defer u.mu.Unlock()

	q := `UPDATE users SET username = $1, display_name = $2 WHERE user_id = $3`
	stmt, err := u.tx.PrepareContext(ctx, q)
	if err != nil {
		return err
	}

	if _, err := stmt.ExecContext(ctx, user.Username, user.DisplayName, user.ID); err != nil {
		return err
	}
	return nil
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
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserPostgres) Search(ctx context.Context, query string) ([]*aggregate.User, error) {
	return nil, nil // this will be implemented using full-text search through ElasticSearch
}

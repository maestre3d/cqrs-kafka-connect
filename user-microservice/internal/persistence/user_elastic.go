package persistence

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/maestre3d/cqrs-kafka-connect/user-microservice/internal/aggregate"
	"github.com/maestre3d/cqrs-kafka-connect/user-microservice/internal/repository"
)

// ReadOnly repository
type UserElastic struct {
	c  *elasticsearch.Client
	mu sync.RWMutex
}

var _ repository.User = &UserElastic{}

func NewUserElastic(c *elasticsearch.Client) *UserElastic {
	return &UserElastic{
		c:  c,
		mu: sync.RWMutex{},
	}
}

func (u *UserElastic) Begin(ctx context.Context) error {
	return nil
}

func (u *UserElastic) Commit() error {
	return nil
}

func (u *UserElastic) Rollback() error {
	return nil
}

func (u *UserElastic) Save(ctx context.Context, user aggregate.User) error {
	return nil
}

func (u *UserElastic) Update(ctx context.Context, user aggregate.User) error {
	return nil
}

func (u *UserElastic) generateId(userID string) string {
	return "Struct{user_id=" + userID + "}"
}

func (u *UserElastic) Find(ctx context.Context, userID string) (*aggregate.User, error) {
	u.mu.RLock()
	defer u.mu.RUnlock()

	res, err := u.c.Get("user.public.users", u.generateId(userID))
	if err != nil {
		return nil, err
	} else if res.StatusCode == 404 {
		return nil, nil
	}
	defer res.Body.Close()

	var e map[string]interface{}
	if err = json.NewDecoder(res.Body).Decode(&e); err != nil {
		return nil, err
	}

	return u.UnmarshalElastic(e["_source"].(map[string]interface{}))
}

func (u *UserElastic) UnmarshalElastic(source map[string]interface{}) (*aggregate.User, error) {
	return &aggregate.User{
		ID:          source["user_id"].(string),
		Username:    source["username"].(string),
		DisplayName: source["display_name"].(string),
	}, nil
}

func (u *UserElastic) Search(ctx context.Context, query string) ([]*aggregate.User, error) {
	return nil, nil
}

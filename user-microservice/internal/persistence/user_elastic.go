package persistence

import (
	"bytes"
	"context"
	"encoding/json"
	"sync"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/maestre3d/cqrs-kafka-connect/user-microservice/internal/aggregate"
	"github.com/maestre3d/cqrs-kafka-connect/user-microservice/internal/domain"
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

	res, err := u.c.Get("user.public.users",
		u.generateId(userID),
		u.c.Get.WithContext(ctx))
	if err != nil {
		return nil, err
	} else if res.StatusCode == 404 {
		res.Body.Close()
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

func (u *UserElastic) Search(ctx context.Context, criteria domain.Criteria) ([]*aggregate.User, error) {
	u.mu.RLock()
	defer u.mu.RUnlock()

	criteria.Fields = []string{"username^2", "*_name"} // set scoring per-field
	criteria.Query = "*" + criteria.Query + "*"        // like statement

	body := newElasticDSLFromCriteria(criteria)
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(body); err != nil {
		return nil, err
	}

	res, err := u.c.Search(
		u.c.Search.WithIndex("user.public.users"),
		u.c.Search.WithContext(ctx),
		u.c.Search.WithBody(&buf),
	)
	if err != nil {
		return nil, err
	} else if res.StatusCode == 404 {
		res.Body.Close()
		return nil, nil
	}
	defer res.Body.Close()

	var e elasticResponseSearch
	if err = json.NewDecoder(res.Body).Decode(&e); err != nil {
		return nil, err
	}

	users := make([]*aggregate.User, 0)
	for _, hit := range e.Hits.Hits {
		user, err := u.UnmarshalElastic(hit.Source)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

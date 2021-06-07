package query

import (
	"context"

	"github.com/maestre3d/cqrs-kafka-connect/user-microservice/internal/aggregate"
	"github.com/maestre3d/cqrs-kafka-connect/user-microservice/internal/application"
	"github.com/maestre3d/cqrs-kafka-connect/user-microservice/internal/domain"
)

type SearchUsers struct {
	Criteria domain.Criteria
}

var _ Query = SearchUsers{}

func (s SearchUsers) Name() string {
	return "user.search"
}

type SearchUsersResponse struct {
	Users         []*aggregate.User `json:"users"`
	TotalSize     int               `json:"total_size"`
	NextPageToken string            `json:"next_page_token"`
}

func newSearchUsersResponse(users []*aggregate.User, nextPage string) SearchUsersResponse {
	return SearchUsersResponse{
		Users:         users,
		TotalSize:     len(users),
		NextPageToken: nextPage,
	}
}

type SearchUsersHandler struct {
	service application.UserReadOnly
}

var _ QueryHandler = &SearchUsersHandler{}

func NewSearchUsersHandler(s application.UserReadOnly) *SearchUsersHandler {
	return &SearchUsersHandler{
		service: s,
	}
}

func (s *SearchUsersHandler) Handle(ctx context.Context, q Query) (interface{}, error) {
	searchQuery := q.(SearchUsers)
	users, nextPage, err := s.service.Search(ctx, searchQuery.Criteria)
	if err != nil {
		return nil, err
	}
	return newSearchUsersResponse(users, nextPage), nil
}

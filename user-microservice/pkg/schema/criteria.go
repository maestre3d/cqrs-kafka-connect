package schema

import (
	"net/http"
	"strconv"

	"github.com/maestre3d/cqrs-kafka-connect/user-microservice/internal/domain"
)

func NewCriteriaFromHTTP(r *http.Request) domain.Criteria {
	pageToken := r.URL.Query().Get("page_token")
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	query := r.URL.Query().Get("query")
	return domain.NewCriteria(pageToken, query, pageSize)
}

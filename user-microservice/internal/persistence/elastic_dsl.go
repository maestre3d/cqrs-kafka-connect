package persistence

import (
	"strconv"

	"github.com/maestre3d/cqrs-kafka-connect/user-microservice/internal/domain"
)

type elasticDSL struct {
	From  int          `json:"from"`
	Size  int          `json:"size"`
	Query elasticQuery `json:"query"`
}

type elasticQuery struct {
	Bool elasticQueryBool `json:"bool"`
}

type elasticQueryBool struct {
	Filters []elasticQueryFilter `json:"filter"`
}

type elasticQueryFilter struct {
	QueryString elasticQueryString `json:"query_string"`
}

type elasticQueryString struct {
	Fields []string `json:"fields"`
	Query  string   `json:"query"`
}

var defaultElasticFromValue = 0

func newElasticDSLFromCriteria(criteria domain.Criteria) elasticDSL {
	pageToken := defaultElasticFromValue
	fromCriteria, err := strconv.Atoi(criteria.PageToken)
	if err == nil {
		pageToken = fromCriteria
	}
	return elasticDSL{
		From: pageToken,
		Size: criteria.PageSize,
		Query: elasticQuery{
			Bool: elasticQueryBool{
				Filters: []elasticQueryFilter{
					{
						QueryString: elasticQueryString{
							Fields: criteria.Fields,
							Query:  criteria.Query,
						},
					},
				},
			},
		},
	}
}

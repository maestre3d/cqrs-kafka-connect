package domain

type Criteria struct {
	PageToken string
	PageSize  int
	Query     string
	Fields    []string
}

var defaultPageSize = 100

func NewCriteria(pageToken, query string, pageSize int, fields ...string) Criteria {
	if pageSize <= 0 || pageSize > 100 {
		pageSize = defaultPageSize
	}
	return Criteria{
		PageToken: pageToken,
		PageSize:  pageSize,
		Query:     query,
		Fields:    fields,
	}
}

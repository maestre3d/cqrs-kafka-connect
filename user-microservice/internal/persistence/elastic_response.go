package persistence

type elasticResponseSearch struct {
	Took     int                  `json:"took"`
	TimedOut bool                 `json:"timed_out"`
	Shards   elasticResponseShard `json:"_shards"`
	Hits     elasticResponseHits  `json:"hits"`
}

type elasticResponseShard struct {
	Total      int `json:"total"`
	Successful int `json:"successful"`
	Skipped    int `json:"skipped"`
	Failed     int `json:"failed"`
}

type elasticResponseHits struct {
	Total    elasticResponseTotal `json:"total"`
	MaxScore float64              `json:"max_score"`
	Hits     []elasticResponseHit `json:"hits"`
}

type elasticResponseTotal struct {
	Value    int    `json:"value"`
	Relation string `json:"relation"`
}

type elasticResponseHit struct {
	Index  string                 `json:"_index"`
	Type   string                 `json:"_type"`
	ID     string                 `json:"_id"`
	Score  float64                `json:"_score"`
	Source map[string]interface{} `json:"_source"`
}

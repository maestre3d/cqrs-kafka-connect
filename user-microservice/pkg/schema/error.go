package schema

import (
	"encoding/json"
	"net/http"
)

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func RespondErrorHTTP(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Add("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(Error{
		Code:    http.StatusInternalServerError,
		Message: err.Error(),
	})
}

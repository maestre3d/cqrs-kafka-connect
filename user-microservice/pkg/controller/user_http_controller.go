package controller

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	nanoid "github.com/matoous/go-nanoid/v2"

	"github.com/maestre3d/cqrs-kafka-connect/user-microservice/internal/application"
	"github.com/maestre3d/cqrs-kafka-connect/user-microservice/pkg/schema"
)

type UserHTTP struct {
	service         application.User
	readOnlyService application.UserReadOnly
}

var _ HTTPController = &UserHTTP{}

func NewUserHTTP(r *mux.Router, s application.User, roService application.UserReadOnly) *UserHTTP {
	u := &UserHTTP{
		service:         s,
		readOnlyService: roService,
	}

	u.MapRoutes(r)
	return u
}

func (c *UserHTTP) MapRoutes(r *mux.Router) {
	r.Path("/users").Methods(http.MethodGet).HandlerFunc(c.list)
	r.Path("/users/{user_id}").Methods(http.MethodGet).HandlerFunc(c.get)
	r.Path("/users").Methods(http.MethodPost).HandlerFunc(c.create)
	r.Path("/users/{user_id}").Methods(http.MethodPut, http.MethodPatch).HandlerFunc(c.update)
}

func (u *UserHTTP) create(w http.ResponseWriter, r *http.Request) {
	userId, err := nanoid.New(16)
	if err != nil {
		schema.RespondErrorHTTP(w, err)
		return
	}
	err = u.service.Create(r.Context(),
		userId,
		r.PostFormValue("username"),
		r.PostFormValue("display_name"))
	if err != nil {
		schema.RespondErrorHTTP(w, err)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		UserID string `json:"user_id"`
	}{
		UserID: userId,
	})
}

func (u *UserHTTP) update(w http.ResponseWriter, r *http.Request) {
	userId := mux.Vars(r)["user_id"]
	err := u.service.Update(r.Context(),
		userId,
		r.PostFormValue("display_name"))
	if err != nil {
		schema.RespondErrorHTTP(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (u *UserHTTP) get(w http.ResponseWriter, r *http.Request) {
	userId := mux.Vars(r)["user_id"]
	user, err := u.readOnlyService.GetById(r.Context(), userId)
	if err != nil {
		schema.RespondErrorHTTP(w, err)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (u *UserHTTP) list(w http.ResponseWriter, r *http.Request) {
	criteria := schema.NewCriteriaFromHTTP(r)
	users, err := u.readOnlyService.Search(r.Context(), criteria)
	if err != nil {
		schema.RespondErrorHTTP(w, err)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

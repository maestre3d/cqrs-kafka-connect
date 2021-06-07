package main

import (
	"database/sql"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"

	"github.com/maestre3d/cqrs-kafka-connect/user-microservice/internal/application"
	"github.com/maestre3d/cqrs-kafka-connect/user-microservice/internal/persistence"
	"github.com/maestre3d/cqrs-kafka-connect/user-microservice/pkg/api"
	"github.com/maestre3d/cqrs-kafka-connect/user-microservice/pkg/controller"
)

func main() {
	esClient, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{"http://127.0.0.1:9200"},
	})
	if err != nil {
		panic(err)
	}

	connStr := "postgres://postgres:root@127.0.0.1:6432/neutrino_users?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	writeRepo := persistence.NewUserPostgres(db)
	readRepo := persistence.NewUserElastic(esClient)
	userService := application.NewUser(writeRepo)
	readService := application.NewUser(readRepo)

	r := mux.NewRouter()
	_ = controller.NewUserHTTP(r, userService, readService)
	srv := api.NewHTTPApi(r)

	panic(srv.ListenAndServe())
}

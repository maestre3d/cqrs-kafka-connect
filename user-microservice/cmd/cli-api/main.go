package main

import (
	"context"
	"database/sql"

	_ "github.com/lib/pq"

	"github.com/maestre3d/cqrs-kafka-connect/user-microservice/internal/application"
	"github.com/maestre3d/cqrs-kafka-connect/user-microservice/internal/persistence"
)

func main() {
	ctx := context.Background()
	connStr := "postgres://postgres:root@127.0.0.1:6432/neutrino_users?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	writeRepo := persistence.NewUserPostgres(db)
	appService := application.NewUser(writeRepo)
	if err := appService.Create(ctx, "1", "br1", "Alfonso Arevalo"); err != nil {
		panic(err)
	}

	if err := appService.Update(ctx, "1", "Bruno Arevalo"); err != nil {
		panic(err)
	}
}

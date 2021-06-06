package application

import "github.com/maestre3d/cqrs-kafka-connect/user-microservice/internal/domain"

func commitOrRollback(tx domain.Transaction, err error) error {
	if err != nil {
		return tx.Rollback()
	}
	return tx.Commit()
}

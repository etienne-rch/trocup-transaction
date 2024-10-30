package services

import (
	"trocup-transaction/repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func DeleteTransaction(id primitive.ObjectID) error {
	// Vérification de l'existence de la transaction avant suppression
	_, err := repository.GetTransactionByID(id)
	if err != nil {
		return err
	}

	// Suppression de la transaction
	return repository.DeleteTransaction(id)
}

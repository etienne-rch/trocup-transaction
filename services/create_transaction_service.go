package services

import (
	"trocup-transaction/models"
	"trocup-transaction/repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateTransaction(transaction *models.Transaction) error {
	// Assigner automatiquement un ID si nécessaire
	if transaction.ID.IsZero() {
		transaction.ID = primitive.NewObjectID()
	}

	// Création de la transaction
	return repository.CreateTransaction(transaction)
}

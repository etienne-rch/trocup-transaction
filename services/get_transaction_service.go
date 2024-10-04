package services

import (
	"trocup-transaction/models"
	"trocup-transaction/repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetTransactionByID(id primitive.ObjectID) (*models.Transaction, error) {
	// Récupération de la transaction par ID
	transaction, err := repository.GetTransactionByID(id)
	if err != nil {
		return nil, err
	}

	// Retourne la transaction récupérée
	return transaction, nil
}
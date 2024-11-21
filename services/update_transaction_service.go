package services

import (
	"errors"
	"trocup-transaction/models"
	"trocup-transaction/repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func UpdateTransaction(id primitive.ObjectID, updatedTransaction *models.Transaction) error {
	// Vérification de l'existence de la transaction à mettre à jour
	existingTransaction, err := repository.GetTransactionByID(id)
	if err != nil {
		return errors.New("transaction not found")
	}

	// Mise à jour des champs modifiables
	existingTransaction.UserB = updatedTransaction.UserB
	existingTransaction.ArticleB = updatedTransaction.ArticleB
	existingTransaction.UserA = updatedTransaction.UserA
	existingTransaction.ArticleA = updatedTransaction.ArticleA
	existingTransaction.Delivery = updatedTransaction.Delivery

	// Enregistrement des modifications dans la base de données
	return repository.UpdateTransaction(id, existingTransaction)
}

func UpdateTransactionState(id primitive.ObjectID, state models.TransactionState) error {
	// Verify transaction exists
	existingTransaction, err := repository.GetTransactionByID(id)
	if err != nil {
		return errors.New("transaction not found")
	}

	// Only update the state field
	existingTransaction.State = state

	// Save to database
	return repository.UpdateTransactionState(id, state)
}

package services

import (
	"trocup-transaction/models"
	"trocup-transaction/repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateTransaction(transaction *models.Transaction) error {
    return repository.CreateTransaction(transaction)
}

func GetTransactions() ([]models.Transaction, error) {
    return repository.GetTransactions()
}

func GetTransactionByID(id primitive.ObjectID) (models.Transaction, error) {
    return repository.GetTransactionByID(id)
}

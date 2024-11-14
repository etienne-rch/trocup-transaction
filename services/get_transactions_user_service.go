package services

import (
	"trocup-transaction/models"
	"trocup-transaction/repository"
)

func GetAllTransactionsByUserID(id string) ([]*models.Transaction, error) {
	// Récupération de toutes les transactions
	transactions, err := repository.GetAllTransactionsByUserID(id)
	if err != nil {
		return nil, err
	}

	// Retourne toutes les transactions récupérées
	return transactions, nil
}

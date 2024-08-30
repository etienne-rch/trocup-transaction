package services

import (
	"trocup-transaction/models"
	"trocup-transaction/repository"
)

func GetAllTransactions() ([]*models.Transaction, error) {
	// Récupération de toutes les transactions
	transactions, err := repository.GetAllTransactions()
	if err != nil {
		return nil, err
	}

	// Retourne toutes les transactions récupérées
	return transactions, nil
}

package handlers

import (
	"net/http"
	"trocup-transaction/services"

	"github.com/gofiber/fiber/v2"
)

func GetAllTransactionsUser(c *fiber.Ctx) error {
	// Récupérer l'ID utilisateur à partir des paramètres
	id := c.Params("id")

	// Appeler le service pour récupérer les transactions
	transactions, err := services.GetAllTransactionsByUserID(id)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve transactions"})
	}

	// Retourner les transactions récupérées
	return c.JSON(transactions)
}

package handlers

import (
	"net/http"
	"trocup-transaction/services"

	"github.com/gofiber/fiber/v2"
)

func GetAllTransactions(c *fiber.Ctx) error {
	transactions, err := services.GetAllTransactions()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve transactions"})
	}

	return c.JSON(transactions)
}

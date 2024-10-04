package handlers

import (
	"net/http"
	"trocup-transaction/services"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetTransaction(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid transaction ID"})
	}

	// Récupérer la transaction
	transaction, err := services.GetTransactionByID(id)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Transaction not found"})
	}

	// Vérifier que l'utilisateur est bien impliqué dans la transaction
	clerkUserId := c.Locals("clerkUserId").(string)
	if transaction.Sender != clerkUserId && transaction.Receiver != clerkUserId {
		return c.Status(http.StatusForbidden).JSON(fiber.Map{"error": "You do not have access to this transaction"})
	}

	return c.JSON(transaction)
}

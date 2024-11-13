package handlers

import (
	"net/http"
	"trocup-transaction/models"
	"trocup-transaction/services"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func UpdateTransaction(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid transaction ID"})
	}

	// Récupérer la transaction pour vérifier le propriétaire
	transaction, err := services.GetTransactionByID(id)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Transaction not found"})
	}

	// Vérifier si l'utilisateur est impliqué dans la transaction
	clerkUserId := c.Locals("clerkUserId").(string)
	if transaction.UserA != clerkUserId && transaction.UserB != clerkUserId {
		return c.Status(http.StatusForbidden).JSON(fiber.Map{"error": "You do not have access to update this transaction"})
	}

	// Mise à jour de la transaction
	var updatedTransaction models.Transaction
	if err := c.BodyParser(&updatedTransaction); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := services.UpdateTransaction(id, &updatedTransaction); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update transaction"})
	}

	return c.JSON(updatedTransaction)
}

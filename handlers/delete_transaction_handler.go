package handlers

import (
	"net/http"
	"trocup-transaction/services"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func DeleteTransaction(c *fiber.Ctx) error {
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

	// Vérifier que l'utilisateur est bien impliqué dans la transaction
	clerkUserId := c.Locals("clerkUserId").(string)
	if transaction.UserA != clerkUserId && transaction.UserB != clerkUserId {
		return c.Status(http.StatusForbidden).JSON(fiber.Map{"error": "You do not have access to delete this transaction"})
	}

	// Supprimer la transaction
	if err := services.DeleteTransaction(id); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete transaction"})
	}

	return c.SendStatus(http.StatusNoContent)
}

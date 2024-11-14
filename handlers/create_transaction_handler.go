package handlers

import (
	"net/http"
	"trocup-transaction/models"
	"trocup-transaction/services"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validate = validator.New()

func CreateTransaction(c *fiber.Ctx) error {
	var transaction models.Transaction
	if err := c.BodyParser(&transaction); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Récupérer l'utilisateur authentifié
	clerkUserId := c.Locals("clerkUserId").(string)

	// Vérifier si l'utilisateur est bien le vendeur (expéditeur)
	if transaction.UserA != clerkUserId {
		return c.Status(http.StatusForbidden).JSON(fiber.Map{"error": "You do not have permission to create this transaction"})
	}

	// Validation de la transaction
	if err := validate.Struct(transaction); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// Créer la transaction dans la base de données
	if err := services.CreateTransaction(&transaction); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create transaction"})
	}

	return c.Status(http.StatusCreated).JSON(transaction)
}

package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"trocup-transaction/config"
	"trocup-transaction/handlers"
	"trocup-transaction/models"
	"trocup-transaction/repository"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestUpdateTransaction(t *testing.T) {
	app := fiber.New()

	// Créer une transaction pour les tests
	id := primitive.NewObjectID()
	transaction := models.Transaction{
		ID:       id,
		UserA:    "receiverUserId456", // Utilise un string pour Receiver
		ArticleB: primitive.NewObjectID(),
		UserB:    "senderUserId123", // Utilise un string pour Sender
		Delivery: &models.Delivery{
			PackageWeight: 2,
			Cost:          100,
			QrCodeUrl:     "http://example.com/qrcode",
		},
	}

	// Créer la transaction dans la base de données
	err := repository.CreateTransaction(&transaction)
	if err != nil {
		t.Fatalf("Failed to create transaction: %v", err)
	}

	// Simuler l'ajout du middleware ClerkAuthMiddleware
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("clerkUserId", transaction.UserB) // Simuler l'utilisateur connecté avec un string (Sender)
		return c.Next()
	})

	// Ajouter le handler de mise à jour de transaction
	app.Put("/transactions/:id", handlers.UpdateTransaction)

	// Modifier la transaction pour la mise à jour
	transaction.Delivery.Cost = 150
	reqBody, _ := json.Marshal(transaction)

	// Créer la requête PUT
	req := httptest.NewRequest("PUT", "/transactions/"+id.Hex(), bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	// Exécuter la requête PUT
	resp, _ := app.Test(req)

	// Vérifier le code de statut attendu
	utils.AssertEqual(t, http.StatusOK, resp.StatusCode, "Expected status code to be 200 OK")

	// Nettoyage après chaque test
	defer config.CleanUpTestDatabase("test_db")
}

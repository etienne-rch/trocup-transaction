package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"trocup-transaction/config"
	"trocup-transaction/handlers"
	"trocup-transaction/models"
	"trocup-transaction/repository"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestGetTransaction(t *testing.T) {
	app := fiber.New()

	// Créer une transaction pour les tests
	id := primitive.NewObjectID() // ID de la transaction reste un ObjectID pour MongoDB
	transaction := models.Transaction{
		ID:            id,
		Receiver:      "receiverUserId456",     // Utilise un string pour Receiver
		SenderArticle: primitive.NewObjectID(), // Article reste un ObjectID
		Sender:        "senderUserId123",       // Utilise un string pour Sender
		Delivery: models.Delivery{
			Type:          "standard",
			PackageWeight: 2,
			Sent:          time.Now(),
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
		c.Locals("clerkUserId", transaction.Sender) // Simuler l'utilisateur connecté avec un string
		return c.Next()
	})

	// Ajouter le handler de récupération de transaction
	app.Get("/transactions/:id", handlers.GetTransaction)

	// Créer la requête GET
	req := httptest.NewRequest("GET", "/transactions/"+id.Hex(), nil)

	// Exécuter la requête GET
	resp, _ := app.Test(req)

	// Vérifier le code de statut attendu
	utils.AssertEqual(t, http.StatusOK, resp.StatusCode, "Expected status code to be 200 OK")

	// Nettoyage après chaque test
	defer config.CleanUpTestDatabase("test_db")
}

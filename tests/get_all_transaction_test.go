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

func TestGetAllTransactions(t *testing.T) {
	app := fiber.New()

	// Ajouter le handler pour récupérer toutes les transactions
	app.Get("/transactions", handlers.GetAllTransactions)

	// Créer deux transactions avec des strings pour Sender et Receiver
	transaction1 := models.Transaction{
		Receiver: "receiverUserId456", // Utilise un string pour Receiver
		Article:  primitive.NewObjectID(),
		Sender:   "senderUserId123", // Utilise un string pour Sender
		Delivery: models.Delivery{
			Type:          "standard",
			PackageWeight: 2,
			Sent:          time.Now(),
			Cost:          100,
			QrCodeUrl:     "http://example.com/qrcode",
		},
	}
	transaction2 := models.Transaction{
		Receiver: "receiverUserId789", // Utilise un string pour Receiver
		Article:  primitive.NewObjectID(),
		Sender:   "senderUserId456", // Utilise un string pour Sender
		Delivery: models.Delivery{
			Type:          "express",
			PackageWeight: 1,
			Sent:          time.Now(),
			Cost:          200,
			QrCodeUrl:     "http://example.com/qrcode2",
		},
	}

	// Créer les transactions dans la base de données
	_ = repository.CreateTransaction(&transaction1)
	_ = repository.CreateTransaction(&transaction2)

	// Créer la requête GET pour récupérer toutes les transactions
	req := httptest.NewRequest("GET", "/transactions", nil)
	resp, _ := app.Test(req)

	// Vérifier que la requête retourne le bon code de statut
	utils.AssertEqual(t, http.StatusOK, resp.StatusCode, "Expected status code to be 200 OK")

	// Nettoyage après chaque test
	defer config.CleanUpTestDatabase("test_db")
}

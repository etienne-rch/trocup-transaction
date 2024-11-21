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

func TestGetAllTransactionsUser(t *testing.T) {
	app := fiber.New()

	// Ajouter le handler pour récupérer les transactions d'un utilisateur
	app.Get("/transactions/users/:id", handlers.GetAllTransactionsUser)

	// Créer trois transactions avec des strings pour UserA et UserB
	transaction1 := models.Transaction{
		UserA:    "userA123",
		UserB:    "userB123",
		ArticleA: &primitive.ObjectID{},
		ArticleB: primitive.NewObjectID(),
		Delivery: &models.Delivery{
			PackageWeight: 2,
			Cost:          100,
			QrCodeUrl:     "http://example.com/qrcode1",
		},
		CreatedAt: time.Now(),
	}
	transaction2 := models.Transaction{
		UserA:    "userA456",
		UserB:    "userB456",
		ArticleA: &primitive.ObjectID{},
		ArticleB: primitive.NewObjectID(),
		Delivery: &models.Delivery{
			PackageWeight: 1,
			Cost:          200,
			QrCodeUrl:     "http://example.com/qrcode2",
		},
		CreatedAt: time.Now(),
	}
	transaction3 := models.Transaction{
		UserA:    "userA123",
		UserB:    "userB789",
		ArticleA: &primitive.ObjectID{},
		ArticleB: primitive.NewObjectID(),
		Delivery: &models.Delivery{
			PackageWeight: 3,
			Cost:          150,
			QrCodeUrl:     "http://example.com/qrcode3",
		},
		CreatedAt: time.Now(),
	}

	// Insérer les transactions dans la base de données
	_ = repository.CreateTransaction(&transaction1)
	_ = repository.CreateTransaction(&transaction2)
	_ = repository.CreateTransaction(&transaction3)

	// Faire la requête GET pour récupérer les transactions d'un utilisateur particulier
	req := httptest.NewRequest("GET", "/transactions/users/userA123", nil)
	resp, _ := app.Test(req)

	// Vérifier que la requête retourne le bon code de statut
	utils.AssertEqual(t, http.StatusOK, resp.StatusCode, "Expected status code to be 200 OK")

	// Nettoyer après le test
	defer config.CleanUpTestDatabase("test_db")
}

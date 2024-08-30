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

	app.Get("/transactions", handlers.GetAllTransactions)

	transaction1 := models.Transaction{
		Version:  1,
		Receiver: primitive.NewObjectID(),
		Article:  primitive.NewObjectID(),
		Sender:   primitive.NewObjectID(),
		Delivery: models.Delivery{
			Type:          "standard",
			PackageWeight: 2,
			Sent:          time.Now(),
			Cost:          100,
			QrCodeUrl:     "http://example.com/qrcode",
		},
	}
	transaction2 := models.Transaction{
		Version:  1,
		Receiver: primitive.NewObjectID(),
		Article:  primitive.NewObjectID(),
		Sender:   primitive.NewObjectID(),
		Delivery: models.Delivery{
			Type:          "express",
			PackageWeight: 1,
			Sent:          time.Now(),
			Cost:          200,
			QrCodeUrl:     "http://example.com/qrcode2",
		},
	}

	_ = repository.CreateTransaction(&transaction1)
	_ = repository.CreateTransaction(&transaction2)

	req := httptest.NewRequest("GET", "/transactions", nil)
	resp, _ := app.Test(req)

	utils.AssertEqual(t, http.StatusOK, resp.StatusCode, "Expected status code to be 200 OK")

	// Nettoyage après chaque test
	defer config.CleanUpTestDatabase("test_db")
}

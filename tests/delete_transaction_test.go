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

func TestDeleteTransaction(t *testing.T) {
	app := fiber.New()
	
	app.Delete("/transactions/:id", handlers.DeleteTransaction)

	id := primitive.NewObjectID()
	transaction := models.Transaction{
		ID:       id,
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
	_ = repository.CreateTransaction(&transaction)

	req := httptest.NewRequest("DELETE", "/transactions/"+id.Hex(), nil)
	resp, _ := app.Test(req)

	utils.AssertEqual(t, http.StatusNoContent, resp.StatusCode, "Expected status code to be 204 No Content")

	// Nettoyage après chaque test
	defer config.CleanUpTestDatabase("test_db")
}

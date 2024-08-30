package tests

import (
	"bytes"
	"encoding/json"
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

func TestUpdateTransaction(t *testing.T) {
	app := fiber.New()

	app.Put("/transactions/:id", handlers.UpdateTransaction)

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

	transaction.Delivery.Cost = 150
	reqBody, _ := json.Marshal(transaction)
	req := httptest.NewRequest("PUT", "/transactions/"+id.Hex(), bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)

	utils.AssertEqual(t, http.StatusOK, resp.StatusCode, "Expected status code to be 200 OK")

	// Nettoyage après chaque test
	defer config.CleanUpTestDatabase("test_db")
}

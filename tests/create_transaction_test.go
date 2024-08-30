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

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestCreateTransaction(t *testing.T) {
	app := fiber.New()
	
	app.Post("/transactions", handlers.CreateTransaction)

	// Transaction sans les champs non obligatoires (Sent, QrCodeUrl)
	transaction := models.Transaction{
		Version:  1,
		Receiver: primitive.NewObjectID(),
		Article:  primitive.NewObjectID(),
		Sender:   primitive.NewObjectID(),
		Delivery: models.Delivery{
			Type:          "standard",
			PackageWeight: 2,
			Cost:          100,
		},
	}

	reqBody, _ := json.Marshal(transaction)
	req := httptest.NewRequest("POST", "/transactions", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)

	utils.AssertEqual(t, http.StatusCreated, resp.StatusCode, "Expected status code to be 201 Created")

	// Transaction avec un champ manquant pour tester la validation
	invalidTransaction := models.Transaction{
		Version:  1,
		Article:  primitive.NewObjectID(),
		Sender:   primitive.NewObjectID(),
		Delivery: models.Delivery{
			Type:          "standard",
			PackageWeight: 2,
			Cost:          100,
		},
	}

	reqBodyInvalid, _ := json.Marshal(invalidTransaction)
	reqInvalid := httptest.NewRequest("POST", "/transactions", bytes.NewReader(reqBodyInvalid))
	reqInvalid.Header.Set("Content-Type", "application/json")
	respInvalid, _ := app.Test(reqInvalid)

	utils.AssertEqual(t, http.StatusBadRequest, respInvalid.StatusCode, "Expected status code to be 400 Bad Request")

	// Nettoyage après chaque test
	defer config.CleanUpTestDatabase("test_db")
}

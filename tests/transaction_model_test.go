package tests

import (
	"testing"
	"time"
	"trocup-transaction/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestTransaction(t *testing.T) {
	id := "507f1f77bcf86cd799439011" // Utilise un string pour Clerk's user IDs

	now := time.Now()

	delivery := models.Delivery{
		ID:            primitive.NewObjectID(),
		Type:          "standard",
		PackageWeight: 1000,
		Sent:          now,
		Cost:          10,
		QrCodeUrl:     "https://example.com/qrcode",
	}

	transaction := models.Transaction{
		ID:       primitive.NewObjectID(), // L'ID de la transaction reste un ObjectID
		UserB:    id,                      // Utilise un string pour UserB
		ArticleA: primitive.NewObjectID(),
		UserA:    id, // Utilise un string pour UserA
		Delivery: delivery,
	}

	// Tests pour vérifier les champs de la transaction
	if transaction.UserB != id {
		t.Errorf("expected UserB to be %v, got %v", id, transaction.UserB)
	}
	if transaction.UserA != id {
		t.Errorf("expected UserA to be %v, got %v", id, transaction.UserA)
	}
	if transaction.Delivery.ID != delivery.ID {
		t.Errorf("expected Delivery ID to be %v, got %v", delivery.ID, transaction.Delivery.ID)
	}
	if transaction.Delivery.Type != "standard" {
		t.Errorf("expected Delivery Type to be 'standard', got %s", transaction.Delivery.Type)
	}
	if transaction.Delivery.PackageWeight != 1000 {
		t.Errorf("expected Package Weight to be 1000, got %d", transaction.Delivery.PackageWeight)
	}
	if !transaction.Delivery.Sent.Equal(now) {
		t.Errorf("expected Sent to be %v, got %v", now, transaction.Delivery.Sent)
	}
	if transaction.Delivery.Cost != 10 {
		t.Errorf("expected Cost to be 10, got %d", transaction.Delivery.Cost)
	}
	if transaction.Delivery.QrCodeUrl != "https://example.com/qrcode" {
		t.Errorf("expected QrCodeUrl to be 'https://example.com/qrcode', got %s", transaction.Delivery.QrCodeUrl)
	}
}

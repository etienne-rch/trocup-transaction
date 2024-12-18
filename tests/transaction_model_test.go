package tests

import (
	"testing"
	"trocup-transaction/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestTransaction(t *testing.T) {
	id := "507f1f77bcf86cd799439011" // Utilise un string pour Clerk's user IDs

	delivery := models.Delivery{
		ID:            primitive.NewObjectID(),
		PackageWeight: 1000,
		Cost:          10,
		QrCodeUrl:     "https://example.com/qrcode",
	}

	transaction := models.Transaction{
		ID:       primitive.NewObjectID(), // L'ID de la transaction reste un ObjectID
		UserA:    id,                      // Utilise un string pour Receiver
		ArticleB: primitive.NewObjectID(),
		UserB:    id, // Utilise un string pour Sender
		Delivery: &delivery,
	}

	// Tests pour vérifier les champs de la transaction
	if transaction.UserA != id {
		t.Errorf("expected UserA to be %v, got %v", id, transaction.UserA)
	}
	if transaction.UserB != id {
		t.Errorf("expected UserB to be %v, got %v", id, transaction.UserB)
	}
	if transaction.Delivery.ID != delivery.ID {
		t.Errorf("expected Delivery ID to be %v, got %v", delivery.ID, transaction.Delivery.ID)
	}
	if transaction.Delivery.PackageWeight != 1000 {
		t.Errorf("expected Package Weight to be 1000, got %d", transaction.Delivery.PackageWeight)
	}
	if transaction.Delivery.Cost != 10 {
		t.Errorf("expected Cost to be 10, got %d", transaction.Delivery.Cost)
	}
	if transaction.Delivery.QrCodeUrl != "https://example.com/qrcode" {
		t.Errorf("expected QrCodeUrl to be 'https://example.com/qrcode', got %s", transaction.Delivery.QrCodeUrl)
	}
}

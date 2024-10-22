package tests

import (
	"testing"
	"time"
	"trocup-transaction/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestTransaction(t *testing.T) {
	id, err := primitive.ObjectIDFromHex("507f1f77bcf86cd799439011")
	if err != nil {
		t.Fatalf("failed to create ObjectID: %v", err)
	}

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
		ID:       id,
		Version:  1,
		Receiver: id,
		Article:  id,
		Sender:   id,
		Delivery: delivery,
	}

	if transaction.ID != id {
		t.Errorf("expected ID to be %v, got %v", id, transaction.ID)
	}
	if transaction.Version != 1 {
		t.Errorf("expected Version to be 1, got %d", transaction.Version)
	}
	if transaction.Receiver != id {
		t.Errorf("expected Receiver to be %v, got %v", id, transaction.Receiver)
	}
	if transaction.Article != id {
		t.Errorf("expected Article to be %v, got %v", id, transaction.Article)
	}
	if transaction.Sender != id {
		t.Errorf("expected Sender to be %v, got %v", id, transaction.Sender)
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

package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Delivery struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Type          string             `bson:"type" json:"type"`
	PackageWeight int                `bson:"packageWeight" json:"packageWeight"`
	Sent          time.Time          `bson:"sent,omitempty" json:"sent,omitempty"`
	Cost          int                `bson:"cost" json:"cost"`
	QrCodeUrl     string             `bson:"qrCodeUrl,omitempty" json:"qrCodeUrl,omitempty"`
}

type Transaction struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Receiver      string             `bson:"receiver" json:"receiver" validate:"required"`
	SenderArticle primitive.ObjectID `bson:"senderArticle" json:"senderArticle" validate:"required"`
	Sender        string             `bson:"sender" json:"sender" validate:"required"`
	Delivery      Delivery           `bson:"delivery" json:"delivery"`
}

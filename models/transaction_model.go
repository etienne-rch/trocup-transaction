package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Delivery struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Type          string             `bson:"type" json:"type"`
	PackageWeight int                `bson:"packageWeight" json:"packageWeight"`
	SentAt        time.Time          `bson:"sentAt,omitempty" json:"sentAt,omitempty"`
	ReceivedAt    time.Time          `bson:"receivedAt,omitempty" json:"receivedAt,omitempty"`
	Cost          int                `bson:"cost" json:"cost"`
	QrCodeUrl     string             `bson:"qrCodeUrl,omitempty" json:"qrCodeUrl,omitempty"`
}

type Transaction struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	State         string             `bson:"state" json:"state"`
	UserA         string             `bson:"userA" json:"userA" validate:"required"`
	ArticleA      primitive.ObjectID `bson:"articleA" json:"articleA" validate:"required"`
	UserB         string             `bson:"userB" json:"userB" validate:"required"`
	ArticleB      primitive.ObjectID `bson:"articleB,omitempty" json:"articleB,omitempty"`
	Delivery      Delivery           `bson:"delivery" json:"delivery,omitempty"`
	CreatedAt     time.Time          `bson:"createdAt" json:"createdAt"`
}

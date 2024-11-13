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
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	UserA     string             `bson:"userA" json:"userA" validate:"required"`
	ArticleA  primitive.ObjectID `bson:"articleA" json:"articleA" validate:"required"`
	UserB     string             `bson:"userB" json:"userB" validate:"required"`
	ArticleB  primitive.ObjectID `bson:"articleB,omitempty" json:"articleB,omitempty"`
	Delivery  Delivery           `bson:"delivery" json:"delivery"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
}

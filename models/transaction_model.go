package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Delivery struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Type          string             `bson:"type" json:"type" validate:"required,oneof=standard express"`
	PackageWeight int                `bson:"packageWeight" json:"packageWeight" validate:"required,gt=0"`
	Sent          time.Time          `bson:"sent,omitempty" json:"sent,omitempty"`
	Cost          int                `bson:"cost" json:"cost" validate:"required,gt=0"`
	QrCodeUrl     string             `bson:"qrCodeUrl,omitempty" json:"qrCodeUrl,omitempty"`
}

type Transaction struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Receiver string             `bson:"receiver" json:"receiver" validate:"required"`
	Article  primitive.ObjectID `bson:"article" json:"article" validate:"required"`
	Sender   string             `bson:"sender" json:"sender" validate:"required"`
	Delivery Delivery           `bson:"delivery" json:"delivery" validate:"required"`
}

package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TransactionState string

const (
	TransactionStatePending   TransactionState = "PENDING"
	TransactionStateAccepted  TransactionState = "ACCEPTED"
	TransactionStateRefused   TransactionState = "REFUSED"
	TransactionStateCancelled TransactionState = "CANCELLED"
	TransactionStateCompleted TransactionState = "COMPLETED"
)

// Optional: Add validation method
func (s TransactionState) IsValid() bool {
	switch s {
	case TransactionStatePending, TransactionStateAccepted,
		TransactionStateRefused, TransactionStateCancelled:
		return true
	}
	return false
}

type GeoPoints struct {
	Type        string    `json:"type,omitempty" bson:"type,omitempty"`
	Coordinates []float64 `json:"coordinates,omitempty" bson:"coordinates,omitempty"`
}

type Address struct {
	Label     string    `json:"label,omitempty" bson:"label,omitempty"`
	Street    string    `json:"street,omitempty" bson:"street,omitempty"`
	City      string    `json:"city,omitempty" bson:"city,omitempty"`
	Postcode  string    `json:"postcode,omitempty" bson:"postcode,omitempty"`
	Citycode  string    `json:"citycode,omitempty" bson:"citycode,omitempty"`
	Floor     int       `json:"floor,omitempty" bson:"floor,omitempty"`
	Extra     string    `json:"extra,omitempty" bson:"extra,omitempty"`
	GeoPoints GeoPoints `json:"geopoints,omitempty" bson:"geopoints,omitempty"`
}

type Delivery struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Address       Address            `bson:"address" json:"address"`
	Company       string             `bson:"company" json:"company,omitempty"`
	PackageWeight int                `bson:"packageWeight" json:"packageWeight,omitempty"`
	SentAt        time.Time          `bson:"sentAt,omitempty" json:"sentAt,omitempty"`
	ReceivedAt    time.Time          `bson:"receivedAt,omitempty" json:"receivedAt,omitempty"`
	Cost          int                `bson:"cost" json:"cost,omitempty"`
	QrCodeUrl     string             `bson:"qrCodeUrl,omitempty" json:"qrCodeUrl,omitempty"`
}

type Transaction struct {
	ID        primitive.ObjectID  `bson:"_id,omitempty" json:"id,omitempty"`
	State     TransactionState    `bson:"state" json:"state"`
	UserA     string              `bson:"userA" json:"userA"`
	UserB     string              `bson:"userB" json:"userB"`
	ArticleB  primitive.ObjectID  `bson:"articleB" json:"articleB"`
	ArticleA  *primitive.ObjectID `bson:"articleA,omitempty" json:"articleA,omitempty"`
	Delivery  *Delivery           `bson:"delivery,omitempty" json:"delivery,omitempty"`
	CreatedAt time.Time           `bson:"createdAt" json:"createdAt"`
}

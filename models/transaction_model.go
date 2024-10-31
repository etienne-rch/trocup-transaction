package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Delivery struct {
    ID             primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
    Type           string             `bson:"type" json:"type"`
    PackageWeight  int                `bson:"packageWeight" json:"packageWeight"`
    Sent           time.Time          `bson:"sent" json:"sent"`
    Cost           int                `bson:"cost" json:"cost"`
    QrCodeUrl      string             `bson:"qrCodeUrl" json:"qrCodeUrl"`
}

// type Transaction struct {
//     ID         primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
//     Version    int                `bson:"version" json:"version"`
//     Receiver   primitive.ObjectID `bson:"receiver" json:"receiver"`
//     Article    primitive.ObjectID `bson:"article" json:"article"`
//     Sender     primitive.ObjectID `bson:"sender" json:"sender"`
//     Delivery   Delivery           `bson:"delivery" json:"delivery"`
// }

type Transaction struct {
    ID         primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
    Version    int                `bson:"version" json:"version"`
    Receiver   primitive.ObjectID `bson:"receiver" json:"receiver"`
    Article    primitive.ObjectID `bson:"article" json:"article"`
    Sender   primitive.ObjectID `bson:"sender_id" json:"sender_id"`  // Ajout de SenderID
    // Amount     float64            `bson:"amount" json:"amount"`
    Delivery   Delivery           `bson:"delivery" json:"delivery"`
}

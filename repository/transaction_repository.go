package repository

import (
	"context"
	"trocup-transaction/config"
	"trocup-transaction/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var TransactionCollection *mongo.Collection

func InitTransactionRepository() {
    TransactionCollection = config.Client.Database("transaction_dev").Collection("transaction")
}

func CreateTransaction(transaction *models.Transaction) error {
    _, err := TransactionCollection.InsertOne(context.TODO(), transaction)
    return err
}

func GetTransactions() ([]models.Transaction, error) {
    var transactions []models.Transaction
    cursor, err := TransactionCollection.Find(context.TODO(), bson.M{})
    if err != nil {
        return nil, err
    }
    defer cursor.Close(context.TODO())

    for cursor.Next(context.TODO()) {
        var transaction models.Transaction
        if err := cursor.Decode(&transaction); err != nil {
            return nil, err
        }
        transactions = append(transactions, transaction)
    }
    return transactions, nil
}

func GetTransactionByID(id primitive.ObjectID) (models.Transaction, error) {
    var transaction models.Transaction
    filter := bson.M{"_id": id}
    err := TransactionCollection.FindOne(context.TODO(), filter).Decode(&transaction)
    return transaction, err
}

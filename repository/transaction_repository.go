package repository

import (
	"context"
	"trocup-transaction/config"
	"trocup-transaction/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateTransaction(transaction *models.Transaction) error {
	_, err := config.TransactionCollection.InsertOne(context.TODO(), transaction)
	return err
}

func GetTransactionByID(id primitive.ObjectID) (*models.Transaction, error) {
	var transaction models.Transaction
	err := config.TransactionCollection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&transaction)
	return &transaction, err
}

func GetAllTransactions() ([]*models.Transaction, error) {
	var transactions []*models.Transaction
	cursor, err := config.TransactionCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var transaction models.Transaction
		if err := cursor.Decode(&transaction); err != nil {
			return nil, err
		}
		transactions = append(transactions, &transaction)
	}
	return transactions, nil
}

func UpdateTransaction(id primitive.ObjectID, transaction *models.Transaction) error {
	_, err := config.TransactionCollection.UpdateOne(context.TODO(), bson.M{"_id": id}, bson.M{"$set": transaction})
	return err
}

func DeleteTransaction(id primitive.ObjectID) error {
	_, err := config.TransactionCollection.DeleteOne(context.TODO(), bson.M{"_id": id})
	return err
}

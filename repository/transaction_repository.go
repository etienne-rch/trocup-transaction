package repository

import (
	"context"
	"trocup-transaction/config"
	"trocup-transaction/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"bytes"
	"encoding/json"
	"errors"
	"net/http"
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

// UpdateUserBalanceForTransaction met à jour la balance de l'utilisateur en envoyant une requête PUT au microservice utilisateur.
func UpdateUserBalanceForTransaction(userID string, transactionValue float64) error {
	// Préparer les données de la mise à jour
	updateData := map[string]interface{}{
		"transaction_value": transactionValue,
	}
	data, err := json.Marshal(updateData)
	if err != nil {
		return err
	}

	// URL du microservice utilisateur pour la mise à jour de la balance
	userServiceURL := "http://localhost:5001/api/users/" + userID + "/balance"

	// Créer une requête PUT pour mettre à jour la balance
	req, err := http.NewRequest("PUT", userServiceURL, bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	// Effectuer la requête
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Vérifier le statut de la réponse
	if resp.StatusCode != http.StatusOK {
		return errors.New("failed to update user balance")
	}

	return nil
}


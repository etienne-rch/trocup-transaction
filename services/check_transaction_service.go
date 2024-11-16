package services

import (
	"trocup-transaction/repository"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CheckTransactionExists(userA, userB string, articleB primitive.ObjectID, articleA primitive.ObjectID) (bool, error) {
	filter := bson.M{
		"userA": userA,
		"userB": userB,
		"articleB": articleB,
	}
	
	if !articleA.IsZero() {
		filter["articleA"] = articleA
	}
	
	exists, err := repository.TransactionExists(userA, userB, articleB, articleA)
	if err != nil {
		return false, err
	}
	
	return exists, nil
} 
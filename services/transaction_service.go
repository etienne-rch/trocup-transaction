package services

import (
	"trocup-transaction/models"
	"trocup-transaction/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"

)

func CreateTransaction(transaction *models.Transaction) error {
    return repository.CreateTransaction(transaction)
}

func GetTransactions() ([]models.Transaction, error) {
    return repository.GetTransactions()
}

func GetTransactionByID(id primitive.ObjectID) (models.Transaction, error) {
    return repository.GetTransactionByID(id)
}


// Fonction pour mettre à jour le solde de l'utilisateur après une transaction
func UpdateUserBalanceForTransaction(userID string, transactionValue float64) error {
	return repository.UpdateUserBalanceForTransaction(userID, transactionValue)
    // updateData := map[string]interface{}{
    //     "transaction_value": transactionValue,
    // }
    // data, err := json.Marshal(updateData)
    // if err != nil {
    //     return err
    // }

    // // Effectuer la requête PUT pour mettre à jour la balance de l'utilisateur dans le service utilisateur
    // userServiceURL := "http://trocup-user:5001/api/users/" + userID + "/balance"
    // req, err := http.NewRequest("PUT", userServiceURL, bytes.NewBuffer(data))
    // if err != nil {
    //     return err
    // }
    // req.Header.Set("Content-Type", "application/json")

    // client := &http.Client{}
    // resp, err := client.Do(req)
    // if err != nil {
    //     return err
    // }
    // defer resp.Body.Close()

    // if resp.StatusCode != http.StatusOK {
    //     return errors.New("failed to update user balance")
    // }

    // return nil
}
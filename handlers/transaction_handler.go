package handlers

import (
	"net/http"
	"trocup-transaction/models"
	"trocup-transaction/repository"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func SetupRoutes(app *fiber.App) {
    app.Get("/transactions", getTransactions)
    app.Get("/transactions/:id", getTransactionByID)
    app.Post("/transactions", createTransaction)
    app.Post("/transactions", createTransaction)
}

func createTransaction(c *fiber.Ctx) error {
    transaction := new(models.Transaction)
    if err := c.BodyParser(transaction); err != nil {
        return c.Status(http.StatusBadRequest).JSON(fiber.Map{
            "error": err.Error(),
        })
    }

    err := repository.CreateTransaction(transaction)
    if err != nil {
        return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
            "error": err.Error(),
        })
    }
    return c.Status(http.StatusCreated).JSON(transaction)
}

func getTransactions(c *fiber.Ctx) error {
    transactions, err := repository.GetTransactions()
    if err != nil {
        return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
            "error": err.Error(),
        })
    }
    return c.JSON(transactions)
}

func getTransactionByID(c *fiber.Ctx) error {
    idParam := c.Params("id")
    id, err := primitive.ObjectIDFromHex(idParam)
    if err != nil {
        return c.Status(http.StatusBadRequest).JSON(fiber.Map{
            "error": "Invalid ID format",
        })
    }

    transaction, err := repository.GetTransactionByID(id)
    if err != nil {
        return c.Status(http.StatusNotFound).JSON(fiber.Map{
            "error": "Transaction not found",
        })
    }
    return c.JSON(transaction)
}

// // handlers/update_user_handler.go

// func UpdateUserBalanceForTransaction(c *fiber.Ctx) error {
// 	type TransactionRequest struct {
// 		TransactionValue float64 `json:"transaction_value"` // Valeur de l'article acheté
// 	}

// 	var req TransactionRequest
// 	if err := c.BodyParser(&req); err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
// 	}

// 	userID := c.Params("id")
// 	user, err := services.GetUserByID(userID)
// 	if err != nil {
// 		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
// 	}

// 	// Calcul de la nouvelle balance
// 	user.Balance -= req.TransactionValue

// 	// Vérifications pour les limites de balance
// 	if !user.IsPremium && user.Balance < 0 {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Non-premium users cannot have a negative balance"})
// 	}

// 	if user.IsPremium && user.Balance < -user.Credit {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Balance cannot exceed negative credit limit"})
// 	}

// 	updatedUser, err := services.UpdateUser(userID, user)
// 	if err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update user"})
// 	}

// 	return c.JSON(updatedUser)
// }

// func CreateTransactionHandler(c *fiber.Ctx) error {
//     transaction := new(models.Transaction)
//     if err := c.BodyParser(transaction); err != nil {
//         return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
//     }

//     // Appel de la fonction CreateTransaction
//     if err := services.CreateTransaction(transaction); err != nil {
//         return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create transaction"})
//     }

//     // Appel de la fonction UpdateUserBalanceForTransaction pour mettre à jour le solde de l'utilisateur
//     if err := services.UpdateUserBalanceForTransaction(transaction.SenderID, -transaction.Amount); err != nil {
//         return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update user balance"})
//     }

//     return c.Status(http.StatusCreated).JSON(fiber.Map{"message": "Transaction created and user balance updated successfully"})
// }
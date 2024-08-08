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

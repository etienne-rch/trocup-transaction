package routes

import (
	"trocup-transaction/handlers"

	"github.com/gofiber/fiber/v2"
)

func TransactionRoutes(app *fiber.App) {
	app.Get("/health", handlers.HealthCheck)

	api := app.Group("/api")

	api.Post("/transactions", handlers.CreateTransaction)
	api.Get("/transactions/:id", handlers.GetTransaction)
	api.Get("/transactions", handlers.GetAllTransactions)
	api.Put("/transactions/:id", handlers.UpdateTransaction)
	api.Delete("/transactions/:id", handlers.DeleteTransaction)
}

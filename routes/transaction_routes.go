package routes

import (
	"trocup-transaction/handlers"
	"trocup-transaction/middleware"

	"github.com/gofiber/fiber/v2"
)

func TransactionRoutes(app *fiber.App) {
	// Route publique pour vérifier la santé de l'API
	app.Get("/health", handlers.HealthCheck)

	// Groupe des routes de l'API
	api := app.Group("/api")

	// Appliquer le middleware Clerk aux routes sensibles
	api.Use(middleware.ClerkAuthMiddleware)

	api.Post("/transactions", handlers.CreateTransaction)
	api.Get("/transactions/:id", handlers.GetTransaction)
	api.Get("/transactions", handlers.GetAllTransactions)
	api.Put("/transactions/:id", handlers.UpdateTransaction)
	api.Delete("/transactions/:id", handlers.DeleteTransaction)
}

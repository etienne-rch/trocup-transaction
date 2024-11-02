package routes

import (
	"trocup-transaction/handlers"
	"trocup-transaction/middleware"

	"github.com/gofiber/fiber/v2"
)

func TransactionRoutes(app *fiber.App) {

	// Routes publiques : accessibles sans authentification
	public := app.Group("/api")

	public.Get("/health", handlers.HealthCheck)


	// Routes protégées : accessibles uniquement avec authentification
	protected := app.Group("/api/protected", middleware.ClerkAuthMiddleware)

	protected.Post("/transactions", handlers.CreateTransaction)
	protected.Get("/transactions/:id", handlers.GetTransaction)
	protected.Get("/transactions", handlers.GetAllTransactions)
	protected.Put("/transactions/:id", handlers.UpdateTransaction)
	protected.Delete("/transactions/:id", handlers.DeleteTransaction)
}

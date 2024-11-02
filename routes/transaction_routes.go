package routes

import (
	"fmt"
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

	// Ajouter une route catch-all pour le débogage
	app.Use(func(c *fiber.Ctx) error {
		return c.Status(404).SendString(fmt.Sprintf("Route not found: %s", c.Path()))
	})
}

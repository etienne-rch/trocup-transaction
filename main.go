package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"trocup-transaction/config"
	"trocup-transaction/routes"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	// Don't fatal if .env doesn't exist - this is expected in production
	_ = godotenv.Load()

	app := fiber.New(fiber.Config{
		Network: "tcp",
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
		AppName: "Trocup Transaction Microservice",
	})

	// Get allowed origins from environment variable
	allowedOrigins := os.Getenv("ALLOWED_ORIGINS")

	// CORS activation for all routes
	app.Use(cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		ExposeHeaders:    "Content-Length",
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Initialize MongoDB
	config.InitMongo()

	// Initialize Clerk
	clerkKey := os.Getenv("CLERK_SECRET_KEY")
	if clerkKey == "" {
		log.Fatal("CLERK_SECRET_KEY is not set")
	}
	clerk.SetKey(clerkKey)

	// Set up routes
	routes.TransactionRoutes(app)

	// Get port from environment variable
	port := os.Getenv("PORT")
	if port == "" {
		port = "5003" // Default port if not specified
	}

	// Create a channel for shutdown signals
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Start server in a goroutine
	serverShutdown := make(chan struct{})
	go func() {
		// Listen on all interfaces (0.0.0.0) instead of just localhost
		if err := app.Listen(fmt.Sprintf("0.0.0.0:%s", port)); err != nil {
			log.Printf("Server error: %v\n", err)
		}
		close(serverShutdown)
	}()

	// Wait for shutdown signal
	select {
	case <-shutdown:
		log.Println("Shutting down server...")
		
		// Create a context with timeout for shutdown
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Shutdown the Fiber app
		if err := app.ShutdownWithContext(ctx); err != nil {
			log.Printf("Server shutdown error: %v\n", err)
		}

		// Disconnect MongoDB
		if err := config.Client.Disconnect(ctx); err != nil {
			log.Printf("MongoDB disconnect error: %v\n", err)
		}

	case <-serverShutdown:
		log.Println("Server stopped unexpectedly")
	}

	log.Println("Server shutdown complete")
}

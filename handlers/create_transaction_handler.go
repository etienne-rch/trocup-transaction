package handlers

import (
	"net/http"
	"os"
	"time"
	"trocup-transaction/models"
	"trocup-transaction/services"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var validate = validator.New()

// Create a struct to match the exact request body
type TransactionRequest struct {
	UserA         string             `json:"userA" validate:"required"`
	ArticleA      primitive.ObjectID `json:"articleA" validate:"required"`
	ArticlePriceA float64           `json:"articlePriceA,omitempty"`
	UserB         string             `json:"userB" validate:"required"`
	ArticleB      primitive.ObjectID `json:"articleB,omitempty"`
	ArticlePriceB float64           `json:"articlePriceB" validate:"required"`
	Delivery      models.Delivery    `json:"delivery"`
	ClerkToken    string            `json:"-"` // The JWT token, "-" means it won't be included in JSON
}

func CreateTransaction(c *fiber.Ctx) error {
	var request TransactionRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Get the JWT token from the request header
	token := c.Get("Authorization")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "No authorization token provided",
		})
	}

	request.ClerkToken = token

	userServiceBaseURL := os.Getenv("USER_SERVICE_URL")

	if userServiceBaseURL == "" {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "User service base URL not set"})
	}

	// Perform a health check on the user microservice
	userServiceHealth := services.GetUserService(userServiceBaseURL).HealthCheck()
	if !userServiceHealth {
		return c.Status(http.StatusServiceUnavailable).JSON(fiber.Map{"error": "User microservice is unavailable"})
	}

	if err := validate.Struct(request); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// Create transaction model from request
	transaction := models.Transaction{
		UserA:     request.UserA,
		UserB:     request.UserB,
		ArticleB:  request.ArticleB,
		Delivery:  request.Delivery,
		CreatedAt: time.Now(),
	}

	// Only set ArticleA if it's not zero
	if !request.ArticleA.IsZero() {
		transaction.ArticleA = request.ArticleA
	}


	// Save to database
	if err := services.CreateTransaction(&transaction); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create transaction"})
	}


	serviceRequest := services.TransactionForUserRequest{
		UserA:         request.UserA,
		UserB:         request.UserB,
		ArticlePriceA: request.ArticlePriceA,
		ArticlePriceB: request.ArticlePriceB,
		ClerkToken:    request.ClerkToken,
	}

	// Update users data on the user microservice
	if err := services.GetUserService(userServiceBaseURL).UpdateUsersData(serviceRequest); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update user"})
	}

	// Update article state on the article microservice
	articleServiceBaseURL := os.Getenv("ARTICLE_SERVICE_URL")
	if articleServiceBaseURL == "" {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Article service base URL not set"})
	}

	articleIDs := []string{request.ArticleA.String(), request.ArticleB.String()}

	if err := services.GetArticleService(articleServiceBaseURL).UpdateArticlesState(
		articleIDs,
		services.ArticleStatusUnavailable,
		request.ClerkToken,
	); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update article"})
	}

	return c.Status(http.StatusCreated).JSON(transaction)
}

package handlers

import (
	"fmt"
	"log"
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


// Article represents the article information in the request
type Article struct {
	ID    string  `json:"id"`
	Price float64 `json:"price"`
}

// TransactionRequest represents the incoming request body
type TransactionRequest struct {
	UserA    string    `json:"userA"`
	UserB    string    `json:"userB"`
	ArticleB Article `json:"articleB"`
	ArticleA Article `json:"articleA,omitempty"`
	State    string  `json:"state"`
	Address  models.Address  `json:"address,omitempty"`
}

func CreateTransaction(c *fiber.Ctx) error {
	var request TransactionRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	log.Printf("Request Body: UserA: %s, UserB: %s, ArticleA: {ID: %s, Price: %.2f}, ArticleB: {ID: %s, Price: %.2f}, State: %s, Address: %s\n",
		request.UserA, request.UserB, request.ArticleA.ID, request.ArticleA.Price, request.ArticleB.ID, request.ArticleB.Price, request.State, request.Address)
	
	token := c.Get("Authorization")

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
	articleBID, err := primitive.ObjectIDFromHex(request.ArticleB.ID)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ArticleB ID format"})
	}

	// Check for existing transaction
	var articleAID primitive.ObjectID
	if request.ArticleA.ID != "" {
		articleAID, err = primitive.ObjectIDFromHex(request.ArticleA.ID)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ArticleA ID format"})
		}
	}

	exists, err := services.CheckTransactionExists(request.UserA, request.UserB, articleBID, articleAID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to check for existing transaction"})
	}
	if exists {
		return c.Status(http.StatusConflict).JSON(fiber.Map{"error": "Transaction already exists"})
	}

	// Continue with creating the transaction...
	transaction := models.Transaction{
		State:     request.State,
		UserA:     request.UserA,
		UserB:     request.UserB,
		ArticleB:  articleBID,
		Delivery: models.Delivery{
			Address: request.Address,
		},
		CreatedAt: time.Now(),
	}

	

	// Only set ArticleA if it's not empty - this is for 1To1 transactions
	if request.ArticleA.ID != "" {
		transaction.ArticleA = articleAID
	}

	// Using fmt.Printf for basic logging
	fmt.Printf("Request Body: %+v\n", request)

	// Or using log package for better logging
	log.Printf("Request Body: %+v\n", request)

	// Save to transaction database
	if err := services.CreateTransaction(&transaction); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create transaction"})
	}


	serviceRequest := services.TransactionForUserRequest{
		UserA:         request.UserA,
		UserB:         request.UserB,
		ArticlePriceA: request.ArticleA.Price,
		ArticlePriceB: request.ArticleB.Price,
		ClerkToken:    token,
	}

	// Update users data on the user microservice
	if request.State == "ACCEPTED" {
		if err := services.GetUserService(userServiceBaseURL).UpdateUsersData(serviceRequest); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update user"})
		}
	}


	// Update article state on the article microservice - only if the transaction is a 1To1 and pending OR if the transaction is a 1ToM and accepted
	articleServiceBaseURL := os.Getenv("ARTICLE_SERVICE_URL")

	if request.State == "ACCEPTED" || (request.State == "PENDING" && request.ArticleA.ID != "") {
		if articleServiceBaseURL == "" {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Article service base URL not set"})
		}
	

		articleIDs := []string{request.ArticleA.ID, request.ArticleB.ID}

		if err := services.GetArticleService(articleServiceBaseURL).UpdateArticlesState(
			articleIDs,
			services.ArticleStatusUnavailable,
			token,
		); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update article"})
		}
	}

	return c.Status(http.StatusCreated).JSON(transaction)
}

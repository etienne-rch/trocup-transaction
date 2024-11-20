package handlers

import (
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
	UserA    string                  `json:"userA"`
	UserB    string                  `json:"userB"`
	ArticleB Article                 `json:"articleB"`
	ArticleA Article                 `json:"articleA,omitempty"`
	State    models.TransactionState `json:"state" validate:"required"`
	Address  models.Address          `json:"address,omitempty"`
}

func CreatePreTransaction(c *fiber.Ctx) error {
	var request TransactionRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Validate state
	if !request.State.IsValid() {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid transaction state",
			"validStates": []string{
				string(models.TransactionStatePending),
				string(models.TransactionStateAccepted),
				string(models.TransactionStateRefused),
				string(models.TransactionStateCancelled),
			},
		})
	}

	log.Printf("Request Body: UserA: %s, UserB: %s, ArticleB: {ID: %s}, State: %s",
		request.UserA,
		request.UserB,
		request.ArticleB.ID,
		request.State)

	if request.ArticleA.ID != "" {
		log.Printf("Optional ArticleA: {ID: %s, Price: %.2f}",
			request.ArticleA.ID,
			request.ArticleA.Price)
	}

	if request.Address.Street != "" || request.Address.City != "" || request.Address.Label != "" ||
		request.Address.Postcode != "" || request.Address.Citycode != "" ||
		len(request.Address.GeoPoints.Coordinates) > 0 {
		log.Printf("Optional Address: Street: %s, City: %s, Label: %s, Postcode: %s, Citycode: %s, GeoPoints: %v",
			request.Address.Street,
			request.Address.City,
			request.Address.Label,
			request.Address.Postcode,
			request.Address.Citycode,
			request.Address.GeoPoints.Coordinates)
	}

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
		CreatedAt: time.Now(),
	}

	// Check if address has any meaningful data
	if request.Address.Street != "" || request.Address.City != "" || request.Address.Label != "" ||
		request.Address.Postcode != "" || request.Address.Citycode != "" ||
		len(request.Address.GeoPoints.Coordinates) > 0 {
		transaction.Delivery = &models.Delivery{
			Address: request.Address,
		}
	}

	// Only set ArticleA if it's provided, for 1To1 transaction
	if request.ArticleA.ID != "" {
		articleAID, err := primitive.ObjectIDFromHex(request.ArticleA.ID)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ArticleA ID format"})
		}
		transaction.ArticleA = &articleAID
	}

	// Save to transaction database
	if err := services.CreateTransaction(&transaction); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create transaction"})
	}

	// Only update article states if this is a 1-to-1 transaction
	if request.ArticleA.ID != "" {
		articleServiceBaseURL := os.Getenv("ARTICLE_SERVICE_URL")
		if articleServiceBaseURL == "" {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Article service base URL not set"})
		}

		// For 1-to-1 transactions, update both articles' states
		articleIDs := []string{request.ArticleB.ID, request.ArticleA.ID}
		_, err = services.GetArticleService(articleServiceBaseURL).UpdateArticlesState(
			articleIDs,
			services.ArticleStatusUnavailable,
			token,
		)
		if err != nil {
			// Rollback the transaction creation
			if rollbackErr := services.DeleteTransaction(transaction.ID); rollbackErr != nil {
				return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update articles and rollback transaction creation"})
			}

			// Rollback the article states if this is a 1-to-1 transaction
			if request.ArticleA.ID != "" {
				articleIDs := []string{request.ArticleB.ID, request.ArticleA.ID}
				_, rollbackErr := services.GetArticleService(articleServiceBaseURL).UpdateArticlesState(
					articleIDs,
					services.ArticleStatusAvailable,
					token,
				)
				if rollbackErr != nil {
					return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update articles and rollback article states"})
				}
			}
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update articles"})
		}
	}

	return c.Status(http.StatusCreated).JSON(transaction)
}

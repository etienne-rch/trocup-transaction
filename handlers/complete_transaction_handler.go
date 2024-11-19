package handlers

import (
    "net/http"
    "os"
    "trocup-transaction/models"
    "trocup-transaction/services"

    "github.com/gofiber/fiber/v2"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

type CompleteTransactionRequest struct {
    State models.TransactionState `json:"state" validate:"required"`
}

func CompleteTransaction(c *fiber.Ctx) error {
    // Parse and validate request body
    var request CompleteTransactionRequest
    if err := c.BodyParser(&request); err != nil {
        return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
    }

    // Validate state using the IsValid method
    if !request.State.IsValid() {
        return c.Status(http.StatusBadRequest).JSON(fiber.Map{
            "error": "Invalid state",
            "validStates": []string{
                string(models.TransactionStatePending),
                string(models.TransactionStateAccepted),
                string(models.TransactionStateRefused),
                string(models.TransactionStateCancelled),
                string(models.TransactionStateCompleted),
            },
        })
    }

    // Get transaction ID from URL parameter
    transactionID := c.Params("id")
    if transactionID == "" {
        return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Transaction ID is required"})
    }

    // Convert to ObjectID
    objectID, err := primitive.ObjectIDFromHex(transactionID)
    if err != nil {
        return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid transaction ID format"})
    }

    // Get transaction
    transaction, err := services.GetTransactionByID(objectID)
    if err != nil {
        return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Transaction not found"})
    }

    // Get clerkUserId from token/context
    clerkUserId := c.Locals("clerkUserId").(string)
    if clerkUserId == "" {
        return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "User not authenticated"})
    }

    // Check if the authenticated user is UserB (the one who receives the transaction request and validates it)
    if clerkUserId != transaction.UserB {
        return c.Status(http.StatusForbidden).JSON(fiber.Map{
            "error": "Only the receiving user can complete this transaction",
        })
    }

	// Check if the transaction is already completed
	if transaction.State == models.TransactionStateCompleted {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Transaction already completed"})
	}

	// Check if the transaction is already refused
	if transaction.State == models.TransactionStateRefused {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Transaction already refused"})

		// If the transaction is refused, we must update the article state to available if it's a 1To1 transaction
		if transaction.ArticleA.IsZero() {
			articleServiceBaseURL := os.Getenv("ARTICLE_SERVICE_URL")
			if articleServiceBaseURL == "" {
				return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Article service base URL not set"})
			}
			
			token := c.Get("Authorization")
			articleIDs := []string{transaction.ArticleB.Hex()}
			_, err := services.GetArticleService(articleServiceBaseURL).UpdateArticlesState(
				articleIDs, 
				services.ArticleStatusAvailable, 
				token,
			)
			if err != nil {
				return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update article status"})
			}
		}
	}

    // Only proceed with user/article updates if state is ACCEPTED
    if request.State == models.TransactionStateAccepted {
        articleServiceBaseURL := os.Getenv("ARTICLE_SERVICE_URL")
        if articleServiceBaseURL == "" {
            return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Article service base URL not set"})
        }

        token := c.Get("Authorization")

        // Collect article IDs
        var articleIDs []string
        articleIDs = append(articleIDs, transaction.ArticleB.Hex())
        if !transaction.ArticleA.IsZero() {
            articleIDs = append(articleIDs, transaction.ArticleA.Hex())
        }

        // Update articles and get their prices in one call
        articles, err := services.GetArticleService(articleServiceBaseURL).UpdateArticlesState(
            articleIDs,
            services.ArticleStatusUnavailable,
            token,
        )
        if err != nil {
            return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update articles"})
        }

        // Extract prices and prepare user service request
        var articlePriceA, articlePriceB float64
        for _, article := range articles {
            if article.ID == transaction.ArticleB.Hex() {
                articlePriceB = article.Price
            }
            if !transaction.ArticleA.IsZero() && article.ID == transaction.ArticleA.Hex() {
                articlePriceA = article.Price
            }
        }

        // Update user data
        userServiceBaseURL := os.Getenv("USER_SERVICE_URL")
        if userServiceBaseURL == "" {
            return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "User service base URL not set"})
        }

        serviceRequest := services.TransactionForUserRequest{
            UserA:         transaction.UserA,
            UserB:         transaction.UserB,
            ArticleB:      transaction.ArticleB.Hex(),
            ArticlePriceB: articlePriceB,
        }

        // Only add ArticleA data if it exists (1-to-1 transaction)
        if !transaction.ArticleA.IsZero() {
            serviceRequest.ArticleA = transaction.ArticleA.Hex()
            serviceRequest.ArticlePriceA = articlePriceA
        }

        if err := services.GetUserService(userServiceBaseURL).UpdateUsersData(serviceRequest, token); err != nil {
            return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update user data"})
        }
    }

    // Update transaction state
    transaction.State = request.State
    if err := services.UpdateTransaction(transaction.ID, transaction); err != nil {
        return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update transaction"})
    }

    // Return appropriate response
    if request.State == models.TransactionStateCompleted {
        return c.Status(http.StatusOK).JSON(fiber.Map{"message": "Transaction completed successfully", "transaction": transaction})
    } else if request.State == models.TransactionStateCancelled {
        return c.Status(http.StatusOK).JSON(fiber.Map{"message": "Transaction cancelled successfully", "transaction": transaction})
    } else {
        return c.Status(http.StatusOK).JSON(fiber.Map{"message": "Transaction updated successfully", "transaction": transaction})
    }
} 
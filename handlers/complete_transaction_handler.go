package handlers

import (
   
    "net/http"
    "os"
    "trocup-transaction/models"
    "trocup-transaction/services"

    "github.com/gofiber/fiber/v2"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "log"
)

type CompleteTransactionRequest struct {
    State models.TransactionState `json:"state" validate:"required"`
}

func CompleteTransaction(c *fiber.Ctx) error {
  
	articleServiceBaseURL := os.Getenv("ARTICLE_SERVICE_URL")
	if articleServiceBaseURL == "" {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Article service base URL not set"})
	}
    
    var request CompleteTransactionRequest
    if err := c.BodyParser(&request); err != nil {
        log.Printf("Error parsing request body: %v", err)
        return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
    }
    log.Printf("Request state: %s", request.State)

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
        log.Printf("Processing ACCEPTED state transaction")
        
        // Only update articles if this is a 1-to-1 transaction
        var articles []services.ArticleUpdateResponse
        var articlePriceA, articlePriceB float64

      

        token := c.Get("Authorization")
        articleIDs := []string{transaction.ArticleB.Hex()}
        
        // Add ArticleA to the update request if it exists
        if !transaction.ArticleA.IsZero() {
            articleIDs = append(articleIDs, transaction.ArticleA.Hex())
        }

        articles, err = services.GetArticleService(articleServiceBaseURL).UpdateArticlesState(
            articleIDs,
            services.ArticleStatusUnavailable,
            token,
        )
        if err != nil {
            return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update articles"})
        }

        // Extract prices from the response
        for _, article := range articles {
            if article.ID == transaction.ArticleB.Hex() {
                articlePriceB = article.Price
            } else if !transaction.ArticleA.IsZero() && article.ID == transaction.ArticleA.Hex() {
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

        // Only add ArticleA fields if it exists and we have its price
        if !transaction.ArticleA.IsZero() && articlePriceA > 0 {
            articleAHex := transaction.ArticleA.Hex()
            serviceRequest.ArticleA = &articleAHex
            serviceRequest.ArticlePriceA = &articlePriceA
        }

        log.Printf("🔥 Content to user service: %s", serviceRequest)

		// Update user data
        if err = services.GetUserService(userServiceBaseURL).UpdateUsersData(serviceRequest, token); err != nil {

			// We roll back and mark the articles as available if the user data update fails
			articles, err = services.GetArticleService(articleServiceBaseURL).UpdateArticlesState(
				articleIDs,
				services.ArticleStatusAvailable,
				token,
			)

			if err != nil {
				return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update articles back to available"})
			}

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
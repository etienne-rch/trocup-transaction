package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"trocup-transaction/handlers"
	"trocup-transaction/models"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Définir un utilisateur simplifié pour les tests
type TestUser struct {
	ID    string // Utilise un string car Clerk utilise des strings pour les ID
	Role  string // Pour définir les permissions (admin, user)
	Email string // Optionnel, si nécessaire pour les tests
}

// Fonction utilitaire pour générer un JWT valide pour un utilisateur de test
func generateJWTForTestUser(user TestUser) (string, error) {
	// Créer un token JWT avec les informations de l'utilisateur
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   user.ID,                               // Clerk's user ID
		"email": user.Email,                            // Facultatif si nécessaire
		"exp":   time.Now().Add(time.Hour * 24).Unix(), // Expire dans 24 heures
	})

	// Signer le token avec ta clé secrète
	secretKey := "secret-key"
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func TestCreateTransaction(t *testing.T) {
	// Créer un utilisateur de test
	testUser := TestUser{
		ID:    "testUserId123",        // Utilise un string comme ID pour Clerk
		Email: "testuser@example.com", // Facultatif
	}

	// Générer un JWT pour cet utilisateur de test
	token, err := generateJWTForTestUser(testUser)
	if err != nil {
		t.Fatalf("Failed to generate JWT: %v", err)
	}

	// Créer une nouvelle application Fiber
	app := fiber.New()

	// Mocker le middleware ClerkAuthMiddleware
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("clerkUserId", testUser.ID) // Mock le middleware en définissant l'utilisateur de test
		return c.Next()
	})

	// Ajouter le handler de création de transaction
	app.Post("/transactions", handlers.CreateTransaction)

	// Créer une transaction valide avec le même Sender que le testUser
	transaction := models.Transaction{
		UserA:     "receiverUserId456",     // Utilise un string pour le Receiver
		ArticleB:  primitive.NewObjectID(), // SenderArticle reste un ObjectID
		UserB:     testUser.ID,             // Utilise l'ID du testUser (string)
		Delivery: models.Delivery{
			Type:          "standard",
			PackageWeight: 2,
			Cost:          100,
		},
	}

	// Créer la requête de création de transaction
	reqBody, _ := json.Marshal(transaction)
	req := httptest.NewRequest("POST", "/transactions", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	// Tester la requête
	resp, _ := app.Test(req)

	// Vérifier que la transaction a été créée avec succès
	utils.AssertEqual(t, http.StatusCreated, resp.StatusCode, "Expected status code to be 201 Created")
}

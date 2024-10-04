package middleware

import (
	"strings"

	"github.com/clerk/clerk-sdk-go/v2/jwt"
	"github.com/gofiber/fiber/v2"
)

// ClerkAuthMiddleware verifies the Authorization header and Clerk session
func ClerkAuthMiddleware(c *fiber.Ctx) error {
	// Extract the Authorization header from the request
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Missing Authorization header",
		})
	}

	// Extract the Bearer token from the Authorization header
	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == authHeader {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid Authorization format",
		})
	}

	// Verify the session token using Clerk's SDK
	claims, err := jwt.Verify(c.Context(), &jwt.VerifyParams{Token: token})
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid or expired token",
		})
	}

	// Store Clerk user data in Fiber's context for future handlers
	c.Locals("clerkUserId", claims.Subject)

	// Continue to the next handler in the chain
	return c.Next()
}

// handlers/logout.go

package handlers // handlers: package for request handlers

import (
	"github.com/go-playground/validator/v10" // validator: package for validating structs
	"github.com/gofiber/fiber/v2"            // fiber: package for building web APIs
)

// Initialize the validator
// This is a global variable that can be accessed from any function
var validate = validator.New()

// Logout handles the user logout process
func Logout(c *fiber.Ctx) error {
	// Since JWTs are stateless, there's no server-side action needed to "invalidate" the token.
	// You can set a flag in the response that tells the client to clear the token.

	return c.Status(fiber.StatusOK).JSON(fiber.Map{ // Return a success response
		"message":    "Logged out successfully", //	Message for successful logout
		"clearToken": true,                      // Flag to clear the token
	})
}

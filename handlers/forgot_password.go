package handlers

import (
	"context"
	"log"
	"math/rand"
	"myfiberproject/database"
	"myfiberproject/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// GenerateRandomPassword generates a random password of the given length
func GenerateRandomPassword(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}

// ForgotPassword handles the forgot password process
func ForgotPassword(c *fiber.Ctx) error {
	// Fetch database name dynamically
	userCollection := database.GetMongoClient().Database(database.GetDatabaseName()).Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	type request struct {
		Email string `json:"email"`
	}

	var req request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	var user models.User
	err := userCollection.FindOne(ctx, bson.M{"email": req.Email, "status": "approved"}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Email not found or not approved"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database error"})
	}

	// Generate a new password
	newPassword := GenerateRandomPassword(10)

	// Hash the new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to hash password"})
	}

	// Update user's password in the database
	_, err = userCollection.UpdateOne(
		ctx,
		bson.M{"_id": user.ID},
		bson.M{"$set": bson.M{"password": string(hashedPassword)}},
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update password"})
	}

	// Return success message
	successMessage := "Password reset successful. Please contact the administrator to retrieve your new password."
	log.Println(successMessage)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": successMessage})
}

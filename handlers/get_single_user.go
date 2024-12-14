package handlers

import (
	"context"
	"myfiberproject/database"
	"myfiberproject/models"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetSingleUser retrieves a single user by their ID
func GetSingleUser(c *fiber.Ctx) error {
	userID, err := primitive.ObjectIDFromHex(c.Params("id")) // Extract the ID from the URL parameter
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	// Dynamically fetch the database name
	userCollection := database.GetMongoClient().Database(database.GetDatabaseName()).Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user models.User
	err = userCollection.FindOne(ctx, bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
		}
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Error querying user database"})
	}

	user.Password = "" // Omit the password in the response for security

	return c.Status(fiber.StatusOK).JSON(user)
}

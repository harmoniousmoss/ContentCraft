package handlers

import (
	"context"
	"myfiberproject/database"
	"myfiberproject/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

// GetAllUsers fetches all user data
func GetAllUsers(c *fiber.Ctx) error {
	// Dynamically fetch the database name
	userCollection := database.GetMongoClient().Database(database.GetDatabaseName()).Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := userCollection.Find(ctx, bson.M{})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error querying user database"})
	}

	var users []models.User
	if err := cursor.All(ctx, &users); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error decoding user data"})
	}

	// Omit the password in the response
	for i := range users {
		users[i].Password = ""
	}

	// Return the users in the response
	return c.Status(fiber.StatusOK).JSON(users)
}

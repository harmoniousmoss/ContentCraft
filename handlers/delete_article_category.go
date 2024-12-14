package handlers

import (
	"context"
	"myfiberproject/database"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func DeleteArticleCategory(c *fiber.Ctx) error {
	// Parse the ID from the URL parameter
	id, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	// Dynamically fetch the database name
	db := database.GetMongoClient().Database(database.GetDatabaseName())
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Delete the document with the matching ID
	_, err = db.Collection("activity_category").DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete data"})
	}

	// Successfully deleted data
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Data successfully deleted"})
}

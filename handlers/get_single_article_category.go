package handlers

import (
	"context"
	"myfiberproject/database"
	"myfiberproject/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetArticleCtegoryByID(c *fiber.Ctx) error {
	// Extracting the ID from the URL parameters
	articleCategoryID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		// If the ID is not a valid ObjectId, return a Bad Request status
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	// Fetch the database name dynamically
	articleCategoryCollection := database.GetMongoClient().Database(database.GetDatabaseName()).Collection("article_category")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var articleCategory models.ArticleCategory
	// Finding the counterpart by its ID in the database
	if err := articleCategoryCollection.FindOne(ctx, bson.M{"_id": articleCategoryID}).Decode(&articleCategory); err != nil {
		// If no document is found, return a Not Found status
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Data not found"})
	}

	// If data is found, return it with an OK status
	return c.Status(fiber.StatusOK).JSON(articleCategory)
}

package handlers

import (
	"context"
	"myfiberproject/database"
	"myfiberproject/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func GetAllArticleCategory(c *fiber.Ctx) error {
	// Dynamically fetch the database name
	articleCategoryCollection := database.GetMongoClient().Database(database.GetDatabaseName()).Collection("article_category")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := articleCategoryCollection.Find(ctx, bson.M{})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error querying database"})
	}

	var articleCategory []models.ArticleCategory
	if err := cursor.All(ctx, &articleCategory); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error decoding data"})
	}

	return c.Status(fiber.StatusOK).JSON(articleCategory)
}

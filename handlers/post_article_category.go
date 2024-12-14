package handlers

import (
	"myfiberproject/database"
	"myfiberproject/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateArticleCategory(c *fiber.Ctx) error {
	var ArticleCategory models.ArticleCategory

	// Parse the request body into the counterpart struct
	if err := c.BodyParser(&ArticleCategory); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse request"})
	}

	// Validate the struct
	if validationErr := validate.Struct(&ArticleCategory); validationErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": validationErr.Error()})
	}

	// Populate additional fields
	ArticleCategory.ID = primitive.NewObjectID()
	ArticleCategory.CreatedAt = time.Now()
	ArticleCategory.UpdatedAt = time.Now()

	// Fetch the database dynamically
	db := database.GetMongoClient().Database(database.GetDatabaseName())
	collection := db.Collection("article_category")

	// Insert the document into the database
	_, err := collection.InsertOne(c.Context(), ArticleCategory)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to insert data"})
	}

	// Return the created data
	return c.Status(fiber.StatusCreated).JSON(ArticleCategory)
}

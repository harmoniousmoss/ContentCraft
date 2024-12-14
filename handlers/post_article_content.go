package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"myfiberproject/database"
	"myfiberproject/libs"
	"myfiberproject/models"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateArticleContent(c *fiber.Ctx) error {
	// Parse the incoming request body
	var articleContent models.ArticleContent
	if err := c.BodyParser(&articleContent); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse request"})
	}

	// Validate the struct
	if validationErr := validate.Struct(&articleContent); validationErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": validationErr.Error()})
	}

	// Fetch ArticleCategories from the database
	categoryCollection := database.GetMongoClient().Database(database.GetDatabaseName()).Collection("article_category")

	var categories []models.ArticleCategory
	cursor, err := categoryCollection.Find(c.Context(), bson.M{})
	if err != nil {
		log.Println("Failed to fetch article categories:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch article categories"})
	}
	if err := cursor.All(c.Context(), &categories); err != nil {
		log.Println("Failed to parse article categories:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to parse article categories"})
	}

	// Use AI to recommend categories
	recommendedCategories, err := recommendArticleCategories(articleContent.Content, categories)
	if err != nil {
		log.Println("Failed to generate AI recommendations:", err)
		recommendedCategories = []string{}
	}

	// Generate an image using the content
	imageURL, err := generateImageFromContent(articleContent.Content)
	if err != nil {
		log.Println("Failed to generate image:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate image"})
	}

	// Populate fields
	articleContent.ID = primitive.NewObjectID()
	articleContent.CreatedAt = time.Now()
	articleContent.UpdatedAt = time.Now()
	articleContent.Image = imageURL
	articleContent.RecommendedCategories = recommendedCategories

	// Insert the article into the database
	contentCollection := database.GetMongoClient().Database(database.GetDatabaseName()).Collection("article_content")
	_, err = contentCollection.InsertOne(c.Context(), articleContent)
	if err != nil {
		log.Println("Failed to insert article content into the database:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to insert article content"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Article content created successfully",
		"data":    articleContent,
	})
}

func recommendArticleCategories(content string, categories []models.ArticleCategory) ([]string, error) {
	openaiAPIKey := os.Getenv("OPENAI_API_KEY")
	if openaiAPIKey == "" {
		return nil, fmt.Errorf("OpenAI API key not found in environment variables")
	}

	// Prepare the list of category names
	categoryNames := []string{}
	for _, category := range categories {
		categoryNames = append(categoryNames, category.Name)
	}

	// Construct the prompt
	prompt := fmt.Sprintf(`
Based on the article content below, suggest the most relevant categories from the provided list.

Content:
%s

Available Categories (choose relevant ones):
%s

Respond in the following format:
Categories: [comma-separated relevant category names]
`, content, strings.Join(categoryNames, ", "))

	log.Printf("Prompt sent to OpenAI:\n%s", prompt)

	// Call OpenAI API
	response, err := libs.CallOpenAI(openaiAPIKey, prompt)
	if err != nil {
		return nil, err
	}

	log.Printf("Response from OpenAI: %+v", response)

	// Parse AI recommendations
	recommendedCategories := extractRecommendedItems(response.Choices[0].Message.Content, "Categories:")
	return recommendedCategories, nil
}

func generateImageFromContent(content string) (string, error) {
	openaiAPIKey := os.Getenv("OPENAI_API_KEY")
	if openaiAPIKey == "" {
		return "", fmt.Errorf("OpenAI API key not found in environment variables")
	}

	// Truncate the content to ensure it doesn't exceed the API prompt limit
	const maxContentLength = 300
	if len(content) > maxContentLength {
		content = content[:maxContentLength] + "..."
	}

	// Prepare the image generation prompt
	prompt := fmt.Sprintf(
		"Generate a high-quality, visually appealing image based on the following article content. Avoid adding any text or words to the image. Content: %s",
		content,
	)

	// Call the OpenAI Image API
	url := "https://api.openai.com/v1/images/generations"
	payload := map[string]interface{}{
		"prompt": prompt,
		"n":      1,
		"size":   "1024x1024",
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal payload: %v", err)
	}

	// Create the HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", openaiAPIKey))

	// Make the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("OpenAI API error: %s", body)
	}

	// Parse the response
	var result struct {
		Data []struct {
			URL string `json:"url"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %v", err)
	}

	if len(result.Data) > 0 {
		return result.Data[0].URL, nil
	}
	return "", fmt.Errorf("no image URL returned from OpenAI API")
}

// extractRecommendedItems parses a response to extract items based on a given prefix
func extractRecommendedItems(response, prefix string) []string {
	lines := strings.Split(response, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, prefix) {
			// Remove the prefix and clean up the response
			items := strings.TrimPrefix(line, prefix)
			return strings.Split(strings.TrimSpace(items), ", ")
		}
	}
	return []string{}
}

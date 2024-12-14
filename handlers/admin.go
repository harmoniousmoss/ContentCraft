package handlers

import (
	"context"
	"fmt"
	"log"
	"myfiberproject/database"
	"myfiberproject/models"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// SeedAdmin seeds an administrator user
func SeedAdmin() (*models.User, error) {
	// Get MongoDB users collection using the shared GetDatabaseName function
	databaseName := database.GetDatabaseName()
	userCollection := database.GetMongoClient().Database(databaseName).Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Fetch admin details from environment variables
	adminFullName := os.Getenv("ADMIN_SEED_FULLNAME")
	adminEmail := os.Getenv("ADMIN_SEED_EMAIL")
	adminPassword := os.Getenv("ADMIN_SEED_PASSWORD")

	// Set default values for Role and Status
	adminRole := models.Administrator
	adminStatus := models.Approved

	// Validate environment variables
	if adminFullName == "" || adminEmail == "" || adminPassword == "" {
		return nil, fmt.Errorf("one or more required admin seed environment variables are not set")
	}

	// Define the admin user
	admin := &models.User{
		ID:                 primitive.NewObjectID(),
		FullName:           adminFullName,
		Email:              adminEmail,
		Password:           adminPassword, // Will be hashed below
		Role:               adminRole,
		Status:             adminStatus,
		TermsAndConditions: true,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(admin.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("error hashing admin password: %w", err)
	}
	admin.Password = string(hashedPassword)

	// Insert the admin user into the database
	_, err = userCollection.InsertOne(ctx, admin)
	if err != nil {
		return nil, fmt.Errorf("error inserting admin user into the database: %w", err)
	}

	log.Printf("Admin user (%s) seeded successfully!", admin.Email)
	return admin, nil
}

// SeedAdminHandler handles the admin seeding process
func SeedAdminHandler(c *fiber.Ctx) error {
	// Check if the request is authorized
	if c.Get("Authorization") != "Bearer "+os.Getenv("JWT_SECRET") {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	// Call SeedAdmin to perform the actual seeding
	admin, err := SeedAdmin()
	if err != nil {
		log.Printf("Error seeding admin user: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error seeding admin user"})
	}

	// Success response including admin details (excluding password)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":  "Admin user seeded successfully",
		"admin_id": admin.ID.Hex(),
		"email":    admin.Email,
		"role":     admin.Role,
		"status":   admin.Status,
	})
}

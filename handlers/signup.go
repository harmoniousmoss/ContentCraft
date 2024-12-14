package handlers

import (
	"context"
	"log"
	"myfiberproject/database"
	"myfiberproject/models"
	"os"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func Signup(c *fiber.Ctx) error {
	log.Printf("Received request body: %s\n", c.Body())
	databaseName := database.GetDatabaseName()
	userCollection := database.GetMongoClient().Database(databaseName).Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Parse request body into User model
	var user models.User
	if err := c.BodyParser(&user); err != nil {
		log.Printf("Error parsing request body: %v\n", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	// Check if the email already exists
	var existingUser models.User
	err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&existingUser)
	if err != mongo.ErrNoDocuments {
		log.Printf("Email already in use: %s\n", user.Email)
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Email is already in use"})
	}

	// Assign default Role and Status
	user.Role = models.Viewer
	user.Status = models.Pending

	// Validate the user data
	if err = validator.New().Struct(user); err != nil {
		var errors []string
		for _, err := range err.(validator.ValidationErrors) {
			errors = append(errors, err.StructNamespace()+": "+err.Tag())
			log.Println("Validation error:", err.StructNamespace()+": "+err.Tag())
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Validation failed", "details": errors})
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error hashing password"})
	}
	user.Password = string(hashedPassword)

	// Populate additional fields
	user.ID = primitive.NewObjectID()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	// Insert the user into the database
	_, err = userCollection.InsertOne(ctx, user)
	if err != nil {
		log.Printf("Error inserting user into database: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot insert user"})
	}
	log.Printf("User %s successfully created with ID: %s\n", user.Email, user.ID.Hex())

	// Create JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":     user.ID.Hex(),
		"role":   user.Role,
		"status": user.Status,
		"exp":    time.Now().Add(time.Hour * 2).Unix(), // Token expires in 2 hours
	})

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Println("JWT_SECRET is not set in environment variables")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "JWT_SECRET environment variable not set"})
	}

	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		log.Printf("Error signing token: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error signing token"})
	}

	// Omit the password in the response and send the token back to the client
	user.Password = ""
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"user": user, "token": tokenString})
}

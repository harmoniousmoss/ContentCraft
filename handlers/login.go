package handlers

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"regexp"
	"strconv"
	"time"

	"myfiberproject/database"
	"myfiberproject/models"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Email      string `json:"email"`
	Password   string `json:"password"`
	RememberMe *bool  `json:"rememberMe"` // Ensure JSON key matches exactly as in the JSON payload.
}

// SafeLogRequestBody safely logs the request body by masking passwords
func SafeLogRequestBody(body []byte) string {
	re := regexp.MustCompile(`("password":")([^"]*)`)
	return re.ReplaceAllString(string(body), `$1[PROTECTED]`)
}

// UnmarshalJSON custom unmarshaler to help debug the boolean field issue.
func (r *LoginRequest) UnmarshalJSON(data []byte) error {
	type Alias LoginRequest
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(r),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	// Log the RememberMe value after unmarshaling to see the actual value
	rememberMeStatus := "nil"
	if r.RememberMe != nil {
		rememberMeStatus = strconv.FormatBool(*r.RememberMe)
	}
	log.Printf("RememberMe after unmarshal: %s", rememberMeStatus)
	return nil
}

func Login(c *fiber.Ctx) error {
	// Fetch database dynamically
	userCollection := database.GetMongoClient().Database(database.GetDatabaseName()).Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	bodyBytes := c.Body()
	safeLogBody := SafeLogRequestBody(bodyBytes)
	log.Printf("Raw body received: %s", safeLogBody)

	var loginRequest LoginRequest
	if err := json.Unmarshal(bodyBytes, &loginRequest); err != nil {
		log.Printf("Error parsing JSON: %v", err)
		return c.Status(fiber.StatusBadRequest).SendString("JSON parsing error")
	}

	// Improved log to show human-readable status for RememberMe
	rememberMeStatus := "not set"
	if loginRequest.RememberMe != nil {
		rememberMeStatus = strconv.FormatBool(*loginRequest.RememberMe)
	}
	log.Printf("Parsed request: Email: %s, RememberMe: %s", loginRequest.Email, rememberMeStatus)

	var user models.User
	err := userCollection.FindOne(ctx, bson.M{"email": loginRequest.Email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Email not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error querying user database"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Incorrect password"})
	}

	var claims jwt.MapClaims
	if loginRequest.RememberMe != nil && *loginRequest.RememberMe {
		claims = jwt.MapClaims{
			"id":     user.ID.Hex(),
			"role":   user.Role,
			"status": user.Status,
		}
	} else {
		expTime := time.Now().Add(72 * time.Hour) // 3 days expiration
		claims = jwt.MapClaims{
			"id":     user.ID.Hex(),
			"role":   user.Role,
			"status": user.Status,
			"exp":    expTime.Unix(),
		}
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Println("JWT_SECRET is not set in environment variables")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "JWT_SECRET environment variable not set"})
	}

	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		log.Printf("Error signing token: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error signing token"})
	}

	log.Printf("User %s logged in with token; RememberMe is set to: %s", loginRequest.Email, rememberMeStatus)

	user.Password = "" // Omit the password for security
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"user": user, "token": tokenString})
}

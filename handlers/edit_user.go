package handlers

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"log"
	"myfiberproject/database"
	"myfiberproject/models"
	"net/http"
	"reflect"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func generatePassword() (string, error) {
	bytes := make([]byte, 12) // Generates a 12-byte password
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	password := base64.URLEncoding.EncodeToString(bytes) // URL-safe, approximately 16 characters
	return password, nil
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return "", err
	}
	return string(hashedPassword), nil
}

func prepareUpdateData(userUpdate models.User) (bson.M, error) {
	updateData := bson.M{"$set": bson.M{}}
	val := reflect.ValueOf(userUpdate)
	for i := 0; i < val.Type().NumField(); i++ {
		field := val.Type().Field(i)
		value := val.Field(i)
		jsonTag := field.Tag.Get("json")

		if jsonTag != "" && jsonTag != "-" && !value.IsZero() {
			if jsonTag == "password" {
				hashedPassword, err := hashPassword(value.String())
				if err != nil {
					return nil, err
				}
				updateData["$set"].(bson.M)[jsonTag] = hashedPassword
			} else {
				updateData["$set"].(bson.M)[jsonTag] = value.Interface()
			}
		}
	}
	updateData["$set"].(bson.M)["updated_at"] = time.Now()
	return updateData, nil
}

func UpdateUser(c *fiber.Ctx) error {
	userID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	var userUpdate models.User
	if err := c.BodyParser(&userUpdate); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Error parsing request body"})
	}

	userCollection := database.GetMongoClient().Database(database.GetDatabaseName()).Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var existingUser models.User
	err = userCollection.FindOne(ctx, bson.M{"_id": userID}).Decode(&existingUser)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Error fetching user from database"})
	}

	updateData, err := prepareUpdateData(userUpdate)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error processing update data"})
	}

	if existingUser.Status == "pending" && userUpdate.Status == "approved" {
		defaultPassword, err := generatePassword()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate password"})
		}
		hashedPassword, err := hashPassword(defaultPassword)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to hash password"})
		}

		updateData["$set"].(bson.M)["password"] = hashedPassword
	}

	result, err := userCollection.UpdateOne(ctx, bson.M{"_id": userID}, updateData)
	if err != nil {
		log.Printf("Error updating user in database: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Error updating user in database"})
	}

	if result.MatchedCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "User updated successfully"})
}

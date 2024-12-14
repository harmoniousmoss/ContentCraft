package middleware

import (
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

const InvalidTokenError = "Invalid token"

type CustomClaims struct {
	Role   string `json:"role"`
	Status string `json:"status"`
	jwt.RegisteredClaims
}

func RequireRole(requiredRoles []string, requiredStatus string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if isPublicAccess(requiredRoles, requiredStatus) {
			return c.Next()
		}

		tokenString, err := getTokenFromHeader(c)
		if err != nil {
			return err
		}

		claims, err := parseToken(c, tokenString)
		if err != nil {
			return err
		}

		// Check if claims is nil before passing it to isRoleAllowed
		if claims == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": InvalidTokenError})
		}

		if isRoleAllowed(claims, requiredRoles, requiredStatus) {
			return c.Next()
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Insufficient permissions"})
	}
}

func isPublicAccess(requiredRoles []string, requiredStatus string) bool {
	return len(requiredRoles) == 1 && requiredRoles[0] == "public" && requiredStatus == ""
}

func getTokenFromHeader(c *fiber.Ctx) (string, error) {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return "", c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "No authorization header provided"})
	}
	return strings.TrimPrefix(authHeader, "Bearer "), nil
}

func parseToken(c *fiber.Ctx, tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		return nil, c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": InvalidTokenError})
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		return nil, c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": InvalidTokenError})
	}
	return claims, nil
}

func isRoleAllowed(claims *CustomClaims, requiredRoles []string, requiredStatus string) bool {
	for _, role := range requiredRoles {
		if claims.Role == role && claims.Status == requiredStatus {
			return true
		}
	}
	return false
}

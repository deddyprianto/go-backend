package middleware

import (
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

func getJWTSecret() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "your-secret-key-here"
	}
	return []byte(secret)
}

func JWTMiddleWare() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		
		if authHeader == "" {
			return c.Status(401).JSON(fiber.Map{
				"error": "Unauthorized",
			})
		}
		// Cek apakah header dimulai dengan "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(401).JSON(fiber.Map{
				"error": "invalid Authorization header format",
			})
		}
		// Ambil token dengan menghapus "Bearer " di awal
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == ""{
			return c.Status(401).JSON(fiber.Map{
				"error": "invalid Authorization header format from TrimPrefix",
			})
		}
		// tokenString := ""
		// if len(authHeader) > 7 && authHeader[:7] == "Bearer" {
		// 	tokenString = authHeader[7:]
		// } else {
		// 	return c.Status(401).JSON(fiber.Map{
		// 		"message": "invalid Authorization header format",
		// 	})
		// }
		
		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return getJWTSecret(), nil
		})

		if err != nil || !token.Valid {
			return c.Status(401).JSON(fiber.Map{
				"message": "sorry your token is expired",
			})
		}

		if claims, ok := token.Claims.(*Claims); ok {
			c.Locals("user_id", claims.UserID)
			c.Locals("user_email", claims.Email)
		}
		return c.Next()
	}
}

func GenerateToken(userId uint, email string) (string, error) {
	claims := Claims{
		UserID: userId,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 1)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(getJWTSecret())
}

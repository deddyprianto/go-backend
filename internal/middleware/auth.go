package middleware

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
type Claims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	Type   string `json:"type"`
	jwt.RegisteredClaims
}

func getJWTSecret() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "your-secret-key-here"
	}
	return []byte(secret)
}
func generateRefreshToken(userId uint, email string) (string, error) {
	claims := Claims{
		UserID: userId,
		Email:  email,
		Type:   "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 30)), // 30 hari
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   "refresh_token",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	fmt.Println("token baru", token)
	return token.SignedString(getJWTSecret())
}
func generateAccessToken(userId uint, email string) (string, error) {
	claims := Claims{
		UserID: userId,
		Email:  email,
		Type:   "access",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)), // 1 hari
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   "access_token",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	fmt.Println("token", token)
	return token.SignedString(getJWTSecret())
}
func GenerateTokenPair(userId uint, email string) (*TokenPair, error) {
	accessToken, err := generateAccessToken(userId, email)
	if err != nil {
		return nil, err
	}
	refreshToken, err := generateRefreshToken(userId, email)
	if err != nil {
		return nil, err
	}
	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
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
		if tokenString == "" {
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

func RefreshAccessToken(refreshTokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(refreshTokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return getJWTSecret(), nil
	})
	if err != nil {
		return "", err
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return "", jwt.ErrInvalidKey
	}

	if claims.Type != "refresh" {
		return "", jwt.NewValidationError("invalid token type", jwt.ValidationErrorClaimsInvalid)
	}
	return generateAccessToken(claims.UserID, claims.Email)
}

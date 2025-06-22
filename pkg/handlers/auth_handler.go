package handlers

import (
	"api-garuda/internal/middleware"
	"api-garuda/pkg/database"
	"api-garuda/pkg/models"
	"database/sql"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	db *sql.DB
}

func NewAuthHandler(db *sql.DB) *AuthHandler {
	return &AuthHandler{
		db: db,
	}
}

func (h *AuthHandler) RegisterUser(c *fiber.Ctx) error {
	var req models.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid request Body: " + err.Error(),
		})
	}
	// validate input
	if req.Email == "" || req.Password == "" || req.Name == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "name or email or password cannot be empty",
		})
	}

	// register user
	response, err := database.Register(h.db, req)
	if err != nil {
		status := 500
		if err.Error() == "user with this email is already exist" {
			status = 409
		}

		return c.Status(status).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(response)
}

func (h *AuthHandler) RefreshTokenHandler(c *fiber.Ctx) error {
	// get secretKey
	secretKey := os.Getenv("JWT_SECRET")
    if secretKey == "" {
        return c.Status(500).JSON(fiber.Map{
            "error": "JWT secret key tidak dikonfigurasi",
        })
    }
    // Ambil token dari header Authorization
    authHeader := c.Get("Authorization")
    if authHeader == "" {
        return c.Status(401).JSON(fiber.Map{
            "error": "Unauthorized, token tidak ditemukan",
        })
    }

    // Extract token dari "Bearer <token>"
    token := strings.TrimPrefix(authHeader, "Bearer ")
    if token == authHeader {
        return c.Status(401).JSON(fiber.Map{
            "error": "Format token tidak valid, gunakan Bearer <token>",
        })
    }

    // Validasi token yang ada
    claims, err := database.ValidateToken(token, secretKey)
    if err != nil {
        return c.Status(401).JSON(fiber.Map{
            "error": "Token tidak valid atau sudah expired",
        })
    }

    // Generate token baru dengan data user yang fresh
    newToken, err := database.GenerateNewToken(claims.UserID, claims.Email, secretKey)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{
            "error": "Gagal generate token baru",
        })
    }

    return c.JSON(fiber.Map{
        "access_token":  newToken.AccessToken,
        "refresh_token": newToken.RefreshToken,
        "user": fiber.Map{
            "id":    claims.UserID,
            "name":  claims.Name,
            "email": claims.Email,
        },
    })
}
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req models.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request Body: " + err.Error(),
		})
	}
	// validate input
	if req.Email == "" || req.Password == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "email or password cannot be empty",
		})
	}

	// login user
	response, err := database.Login(h.db, req)
	if err != nil {
		status := 500
		if err.Error() == "invalid email or password" {
			status = 401
		}
		return c.Status(status).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	tokenPair, err := middleware.GenerateTokenPair(response.User.ID, response.User.Email)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "failed to generate token pair: " + err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"message": "Login success",
		"data": fiber.Map{
			"id":    response.User.ID,
			"email": response.User.Email,
			"name":  response.User.Name,
		},
		"access_token":  tokenPair.AccessToken,
		"refresh_token": tokenPair.RefreshToken,
		"status":        "success",
		"statusCode":    200,
	})

}

func (h *AuthHandler) GetUserProfile(c *fiber.Ctx) error {
    id := c.Locals("user_id").(uint)
    profile, err := database.GetProfile(h.db, id)
    if err != nil {
        return c.Status(404).JSON(fiber.Map{
            "error": err.Error(),
        })
    }
    
    // Pastikan struct response memiliki field yang sesuai
    response := fiber.Map{
		"data": fiber.Map{
		"id":        profile.ID,
        "name":      profile.Name,
        "email":     profile.Email,
        "created_at": profile.CreatedAt,
        "updated_at": profile.UpdatedAt,
		},
		"message": "User profile retrieved successfully",
		"status" : "success",
		"statusCode": 200,
	}

    return c.JSON(response)
}

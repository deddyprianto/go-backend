package handlers

import (
	"api-garuda/pkg/database"
	"api-garuda/pkg/models"
	"database/sql"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	db *sql.DB
}
 
func NewAuthHandler(db *sql.DB) *AuthHandler{
	return &AuthHandler{
		db: db,
	}
}

func (h *AuthHandler) RegisterUser(c *fiber.Ctx) error{
	var req models.RegisterRequest
	fmt.Println(req)
	if err := c.BodyParser(&req); err != nil{
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid request Body: " + err.Error(),
		})
	}
	fmt.Println(req)

	// validate input
	if req.Email != "" || req.Password != "" || req.Name != ""{
		return c.Status(400).JSON(fiber.Map{
			"error" : "name or email or password cannot be empty",
		})
	}

	// register user
	response, err := database.Register(h.db, req)
	if err != nil{
		status := 500
		if err .Error() == "user with this email is already exist"{
			status = 409
		}

		return c.Status(status).JSON(fiber.Map{
			"error" : err.Error(),
		})
	}
	return c.Status(200).JSON(response)
}

func (h *AuthHandler) Login(c *fiber.Ctx) error{
	var req models.LoginRequest
	if err := c.BodyParser(&req); err != nil{
		return c.Status(400).JSON(fiber.Map{
			"error" : "Invalid request Body: " + err.Error(),
		})
	}
	// validate input

	if req.Email == "" || req.Password == ""{
		return c.Status(400).JSON(fiber.Map{
			"error" : "email or password cannot be empty",
		})
	}

	// login user
	response, err := database.Login(h.db,req);

	if err != nil{
		status := 500
		if err.Error() == "invalid email or password"{
			status = 401
		}
		return c.Status(status).JSON(fiber.Map{
			"error" : err.Error(),
		})
	}
	return c.JSON(response)
	
}

func (h *AuthHandler) GetUserProfile(c *fiber.Ctx) error{
	id := c.Locals("user_id").(uint)
	user, err := database.GetProfile(h.db, id)

	if err != nil{
		return c.Status(404).JSON(fiber.Map{
			"error" : err.Error(),
		})
	}

	return c.JSON(user)
}
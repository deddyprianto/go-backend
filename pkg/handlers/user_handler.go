package handlers

import (
	"api-garuda/pkg/database"
	"api-garuda/pkg/models"
	"database/sql"
	"encoding/base64"
	"fmt"
	"io"
	"time"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	db *sql.DB
}

func NewUserHandler(db *sql.DB) *UserHandler {
	return &UserHandler{
		db: db,
	}
}

func (h *UserHandler) GetAllUsers(c *fiber.Ctx) error {
	users, err := database.GetAllUSers(h.db)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Failed to get users: " + err.Error(),
		})
	}
	return c.JSON(users)
}

func (h *UserHandler) GetUserById(c *fiber.Ctx) error {
	id := c.Params("id")
	user, err := database.GetUserById(h.db, id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return c.JSON(user)
}
func (h *UserHandler) CreateEmployee(c *fiber.Ctx) error {
	// ambil form data
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Failed to parse form data: " + err.Error(),
		})
	}
	fmt.Println(form, "form data")
	// ambil nilai dari form fields
	name := form.Value["name"][0]
	position := form.Value["position"][0]

	// ambil dari profile_picture
	file, err := c.FormFile("profile_picture")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "No profile picture uploaded: " + err.Error(),
		})
	}
	// baca file
	fileBuffer, err := file.Open()
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Failed to open file: " + err.Error(),
		})
	}
	defer fileBuffer.Close()
	// baca bytes dari file
	buffer, err := io.ReadAll(fileBuffer)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Failed to read file bytes: " + err.Error(),
		})
	}

	// encode bytes ke base64 strings
	profilePicture := base64.StdEncoding.EncodeToString(buffer)
	// save user to database
	user := models.Employee{
		Name:           name,
		Position:       position,
		CreatedAt:      time.Now(),
		ProfilePicture: profilePicture,
	}
	response, err := database.CreateEmployee(h.db, user)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Failed to create employee: " + err.Error(),
		})
	}

	return c.JSON(response)

}

func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	var user models.User

	if err := c.BodyParser(&user); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Failed to parse user data: " + err.Error(),
		})
	}
	fmt.Println(user, "user data")
	response, err := database.CreateUser(h.db, user)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Failed to create user: " + err.Error(),
		})
	}

	return c.JSON(response)
}

func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	// get ID from params
	id := c.Params("id")

	var user models.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Failed to parse user data: " + err.Error(),
		})
	}
	response, err := database.UpdateUser(h.db, user, id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Failed to update user: " + err.Error(),
		})
	}

	return c.JSON(response)
}

func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	rowsAffected, err := database.DeleteUser(h.db, id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Failed to delete user: " + err.Error(),
		})
	}

	if rowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{
			"message": "User not found",
		})
	}

	return c.JSON(fiber.Map{
		"message": "success delete user",
	})

}

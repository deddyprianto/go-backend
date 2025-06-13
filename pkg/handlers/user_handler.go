package handlers

import (
	"api-garuda/pkg/database"
	"api-garuda/pkg/models"
	"database/sql"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct{
	db *sql.DB
}

func NewUserHandler(db *sql.DB) *UserHandler{
	return &UserHandler{
		db: db,
	}
}

func (h *UserHandler) GetAllUsers(c *fiber.Ctx) error {
	users ,err := database.GetAllUSers(h.db)
	if err != nil{
		return c.Status(500).JSON(fiber.Map{
			"message": "Failed to get users: " + err.Error(),
		})
	}
	return c.JSON(users)
}

func (h *UserHandler) GetUserById(c *fiber.Ctx) error {
	id := c.Params("id")
	user ,err := database.GetUserById(h.db, id)
	if err != nil{
		return c.Status(404).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return c.JSON(user)
}

func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
    var user models.User
	
    if err := c.BodyParser(&user); err != nil {
        return c.Status(400).JSON(fiber.Map{
            "message": "Failed to parse user data: " + err.Error(),
        })
    }    
    id, err := database.CreateUser(h.db, user)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{
            "message": "Failed to create user: " + err.Error(),
        })
    }
    
    return c.JSON(fiber.Map{
        "id": id,
    })
}

func (h *UserHandler) UpdateUser(c *fiber.Ctx) error{
	// get ID from params
	id := c.Params("id");

	var user models.User
	if err := c.BodyParser(&user); err != nil{
		return c.Status(400).JSON(fiber.Map{	
			"message": "Failed to parse user data: " + err.Error(),
		})
	}
	updatedId, err := database.UpdateUser(h.db, user, id)
	if err != nil{
		return c.Status(500).JSON(fiber.Map{
			"message": "Failed to update user: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "success update user",
		"id": updatedId,
	})
}

func (h *UserHandler) DeleteUser(c *fiber.Ctx) error{
	id := c.Params("id")
	rowsAffected, err := database.DeleteUser(h.db, id);
	if err != nil{
		return c.Status(500).JSON(fiber.Map{
			"message": "Failed to delete user: " + err.Error(),
		})
	}

	if rowsAffected == 0{
		return c.Status(404).JSON(fiber.Map{
			"message": "User not found",
		})
	}

	return c.JSON(fiber.Map{
		"message": "success delete user",
	})

}

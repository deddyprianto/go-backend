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
	users , err := database.GetAllUSers(h.db)
	if err != nil{
		return c.Status(500).JSON(fiber.Map{
			"message": err.Error(),
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
            "message": err.Error(),
        })
    }
    
    id, err := database.CreateUser(h.db, user)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{
            "message": err.Error(),
        })
    }
    
    return c.JSON(fiber.Map{
        "id": id,
    })
}

func (h *UserHandler) UpdateUser(c *fiber.Ctx) error{
	// get ID dari params
	id := c.Params("id");

	// lakukan parsing data
	var user models.User
	if err := c.BodyParser(&user); err != nil{
		return c.Status(400).JSON(fiber.Map{	
			"message": "error di awal:" + err.Error(),
		})
	}
	updatedId, err := database.UpdateUser(h.db, user, id)
	if err != nil{
		return c.Status(500).JSON(fiber.Map{
			"PESAN": "ada error di query:" +err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "user berhasil di update",
		"id": updatedId,
	})
}

func (h *UserHandler) DeleteUser(c *fiber.Ctx) error{
	id := c.Params("id")
	rowsAffected, err := database.DeleteUser(h.db, id);
	if err != nil{
		return c.Status(500).JSON(fiber.Map{
			"message": "error di awal" + err.Error(),
		})
	}

	return c.JSONP(fiber.Map{
		"rows_affected": rowsAffected,
	})

}

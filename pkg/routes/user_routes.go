package routes

import (
	"api-garuda/pkg/handlers"
	"database/sql"

	"github.com/gofiber/fiber/v2"
)


func SetupRoutes(app *fiber.App, db *sql.DB){
	userHandler := handlers.NewUserHandler(db)
	
	app.Get("/api/users", userHandler.GetAllUsers)
	app.Get("/api/users/:id", userHandler.GetUserById)
	app.Post("/api/users", userHandler.CreateUser)
	app.Put("/api/users/:id", userHandler.UpdateUser)
	app.Delete("api/users/:id", userHandler.DeleteUser)

}
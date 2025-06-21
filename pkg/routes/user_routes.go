package routes

import (
	"api-garuda/internal/middleware"
	"api-garuda/pkg/handlers"
	"database/sql"

	"github.com/gofiber/fiber/v2"
)


func SetupRoutes(app *fiber.App, db *sql.DB){
	userHandler := handlers.NewUserHandler(db)
	authHandler := handlers.NewAuthHandler(db)

	api := app.Group("/api")

	auth := api.Group("/auth")
	auth.Post("/register", authHandler.RegisterUser)
	auth.Post("/login", authHandler.Login)
	auth.Post("/refresh", authHandler.RefreshTokenHandler)
	
	users := api.Group("/users")

	users.Use(middleware.JWTMiddleWare())

	users.Get("/profile", authHandler.GetUserProfile)
	users.Get("/", userHandler.GetAllUsers)
	users.Get("/:id", userHandler.GetUserById)
	users.Post("/", userHandler.CreateUser)
	users.Put("/:id", userHandler.UpdateUser)
	users.Delete("/:id", userHandler.DeleteUser)
}
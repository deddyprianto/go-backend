package main

import (
	"api-garuda/pkg/database"
	"api-garuda/pkg/routes"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
    db ,err := database.NewConnection()
    if err != nil{
        fmt.Println("FAILED CONNECTED TO DATABASE")
        return
    }
    defer db.Close()

    err = database.PingDatabase(db)
    if err != nil{
        fmt.Println("gagal")
        return
    }
    fmt.Println("FAILED PING TO DATABASE")

    app := fiber.New(fiber.Config{
        ErrorHandler: func(c *fiber.Ctx, err error) error {
            code := fiber.StatusInternalServerError
            if e, ok := err.(*fiber.Error); ok{
                code = e.Code
            }
            return c.Status(code).JSON(fiber.Map{
                "error": err.Error(),
            })
        },
    })

    // middleware
    app.Use(logger.New())
    app.Use(cors.New(cors.Config{
        AllowOrigins: "*",
        AllowMethods: "GET, POST, PUT, DELETE",
        AllowHeaders: "Origin, Content-Type, Accept, Authorization",
    }))
    
    routes.SetupRoutes(app, db)
    log.Fatal(app.Listen(":3000"))
}
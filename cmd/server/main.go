package main

import (
	"api-garuda/pkg/database"
	"api-garuda/pkg/routes"
	"fmt"

	"github.com/gofiber/fiber/v2"
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

    app := fiber.New()
    routes.SetupRoutes(app, db)
    fmt.Println("Server running on port 3000")
    app.Listen(":3000")
}
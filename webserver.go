package main

import (
	"tarjeta/auth"

	fiber "github.com/gofiber/fiber/v2"
)

func getDatabase() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database")
	}
	return db
}

func main() {
	app := fiber.New()

	db := getDatabase()
	db.AutoMigrate

	app.Get("/api/v0", func(c *fiber.Ctx) error {
		return c.SendString("Hello world!")
	})
	app.Post("/api/v0/auth", auth.Login)
	app.Post("/api/v0/logout", auth.Logout)
	app.Get("/api/v0/whoami", auth.WhoAmI)
	app.Static("/", "./public")

	app.Listen(":3000")
}

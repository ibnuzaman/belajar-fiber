package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"time"
)

func main() {
	app := fiber.New(fiber.Config{
		IdleTimeout:  5 * time.Minute,
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
		Prefork:      true,
	})

	app.Use(func(c *fiber.Ctx) error {
		fmt.Println("Sebelum Middleware")
		err := c.Next()
		fmt.Println("Sesudah Middleware")
		return err
	})
	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.SendString("hello world")
	})

	//if fiber.IsChild() {
	//	fmt.Println("child proces")
	//} else {
	//	fmt.Println("parent process")
	//}

	//app.Get("/test", func(c *fiber.Ctx) error {
	//	return c.SendString("Test Get")
	//})

	err := app.Listen("Localhost:8080")
	if err != nil {
		panic(err)
	}
}

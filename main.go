package main

import (
	"fmt"

	"github.com/ahmadammarm/golang-todo-list/db"
	"github.com/ahmadammarm/golang-todo-list/model"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Warning: .env file not found")
	}

	db, err := db.DatabaseInit()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	app := fiber.New()

	// Get all activities route
	app.Get("/activities", func(c *fiber.Ctx) error {
		activities := []model.Activity{}

		rows, err := db.Query("SELECT * FROM activities")
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var activity model.Activity
			if err := rows.Scan(&activity.ID, &activity.Title, &activity.Category, &activity.Description, &activity.ActivityDate, &activity.Status, &activity.CreatedAt); err != nil {
				return err
			}
			activities = append(activities, activity)
		}

		return c.Status(fiber.StatusOK).JSON(activities)
	})

	app.Listen(":8000")
}

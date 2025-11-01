package main

import (
	"database/sql"
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

	// create an activity route
	app.Post("/activities", func(c *fiber.Ctx) error {
		activity := new(model.Activity)
		if err := c.BodyParser(activity); err != nil {
			return err
		}

		query := `INSERT INTO activities (title, category, description, activity_date, status, created_at)
                  VALUES ($1, $2, $3, $4, $5, NOW()) RETURNING id`
		err := db.QueryRow(query, activity.Title, activity.Category, activity.Description, activity.ActivityDate, activity.Status).Scan(&activity.ID)
		if err != nil {
			return err
		}

		return c.Status(fiber.StatusCreated).JSON(activity)
	})

	// get an activity by id route
	app.Get("/activities/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		var activity model.Activity

		query := `SELECT id, title, category, description, activity_date, status, created_at FROM activities WHERE id = $1`
		err := db.QueryRow(query, id).Scan(&activity.ID, &activity.Title, &activity.Category, &activity.Description, &activity.ActivityDate, &activity.Status, &activity.CreatedAt)
		if err != nil {
			if err == sql.ErrNoRows {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Activity not found"})
			}
			return err
		}

		return c.Status(fiber.StatusOK).JSON(activity)
	})

	// edit an activity rby id route
	app.Put("/activities/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		activity := new(model.Activity)
		if err := c.BodyParser(activity); err != nil {
			return err
		}

		query := `UPDATE activities SET title=$1, category=$2, description=$3, activity_date=$4, status=$5 WHERE id=$6`
		result, err := db.Exec(query, activity.Title, activity.Category, activity.Description, activity.ActivityDate, activity.Status, id)
		if err != nil {
			return err
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return err
		}
		if rowsAffected == 0 {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Activity not found"})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Activity updated successfully"})
	})

	// delete an activity by id route
	app.Delete("/activities/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")

		query := `DELETE FROM activities WHERE id=$1`
		result, err := db.Exec(query, id)
		if err != nil {
			return err
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return err
		}
		if rowsAffected == 0 {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Activity not found"})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Activity deleted successfully"})
	})

	app.Listen(":8000")
}

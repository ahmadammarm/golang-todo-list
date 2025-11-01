package model

import "time"

type Activity struct {
	ID           int       `json:"id"`
	Title        string    `json:"title"`
	Category     string    `json:"category"`
	Description  string    `json:"description"`
	ActivityDate time.Time `json:"activity_date"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
}

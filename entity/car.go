package entity

import (
	"time"
)

type Car struct {
	ID             int64     `json:"id"`
	Make           string    `json:"make"`
	Model          string    `json:"model"`
	Package        string    `json:"package"`
	Color          string    `json:"color"`
	Year           int       `json:"year"`
	Category       string    `json:"category"`
	Mileage        int       `json:"mileage"`
	Price          int       `json:"price"`
	Identification string    `json:"identification"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

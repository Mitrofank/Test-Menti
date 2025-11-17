package models

import "time"

type Car struct {
	ID                  int          `json:"id" db:"id"`
	Make                string       `json:"make" db:"make"`
	Model               string       `json:"model" db:"model"`
	Year                int          `json:"year" db:"year"`
	OwnerID             int          `json:"owner_id" db:"owner_id"`
	PreviousOwnersCount int          `json:"previous_owners_count" db:"previous_owners_count"`
	Currency            CurrencyCode `json:"currency" db:"currency"`
	Price               int          `json:"price" db:"price"`
	Options             []string     `json:"options" db:"options"`
	CreatedAt           time.Time    `json:"created_at" db:"created_at"`
}

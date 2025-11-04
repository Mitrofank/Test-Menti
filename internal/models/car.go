package models

type Car struct {
	ID         int      `json:"id" db:"id"`
	Mark       string   `json:"mark" db:"mark"`
	Model      string   `json:"model" db:"model"`
	OwnerCount int      `json:"owner_count" db:"owner_count"`
	Price      int      `json:"price" db:"price"`
	Currency   string   `json:"currency" db:"currency"`
	Options    []string `json:"options" db:"options"`
}

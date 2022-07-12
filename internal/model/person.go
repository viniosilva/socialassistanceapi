package model

import "time"

type Person struct {
	ID        int
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
	AddressID int
	Name      string
}

package model

import "time"

type Person struct {
	ID        int
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
	FamilyID  int
	Name      string
}

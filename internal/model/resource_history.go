package model

import (
	"time"
)

type ResourceToFamily struct {
	ID         int
	CreatedAt  time.Time
	DeletedAt  time.Time
	ResourceID int
	FamilyID   int
	Quantity   float64
}

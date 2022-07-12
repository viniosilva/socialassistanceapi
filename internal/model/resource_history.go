package model

import (
	"time"
)

type ResourceToAddress struct {
	ID         int
	CreatedAt  time.Time
	DeletedAt  time.Time
	ResourceID int
	AddressID  int
	Quantity   float64
}

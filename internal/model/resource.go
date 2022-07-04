package model

import (
	"time"
)

type Resource struct {
	ID        int
	CreatedAt time.Time
	UpdatedAt time.Time

	Name        string
	Amount      float32
	Measurement string
}

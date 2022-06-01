package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type Resource struct {
	ID        int
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time

	Name        string
	Amount      decimal.Decimal
	Measurement string
}

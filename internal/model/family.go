package model

import "time"

type Family struct {
	ID           int
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time
	Name         string
	Country      string
	State        string
	City         string
	Street       string
	Neighborhood string
	Number       string
	Complement   string
	Zipcode      string
}

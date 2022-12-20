package service

type DonateResourceDonateDto struct {
	ResourceID int     `json:"-"`
	AddressID  int     `json:"address_id" example:"1" binding:"required"`
	Quantity   float64 `json:"quantity" example:"10" binding:"required,gte=0"`
}

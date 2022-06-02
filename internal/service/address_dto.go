package service

type Address struct {
	ID           int    `json:"id" example:"1"`
	CreatedAt    string `json:"created_at" example:"2000-01-01T12:03:00"`
	UpdatedAt    string `json:"updated_at" example:"2000-01-01T12:03:00"`
	DeletedAt    string `json:"deleted_at" example:"2000-01-01T12:03:00"`
	Country      string `json:"country" example:"BR"`
	State        string `json:"state" example:"SP"`
	City         string `json:"city" example:"São Paulo"`
	Neighborhood string `json:"neighborhood" example:"Centro Histórico"`
	Street       string `json:"street" example:"R. Vinte e Cinco de Março"`
	Number       string `json:"number" example:"1000"`
	Complement   string `json:"complement" example:"1A"`
	Zipcode      string `json:"zipcode" example:"01021100"`
}

type AddressResponse struct {
	Data *Address `json:"data"`
}

type AddressesResponse struct {
	Data []Address `json:"data"`
}

type AddressDto struct {
	Country      string `json:"country" example:"BR" binding:"required"`
	State        string `json:"state" example:"SP" binding:"required"`
	City         string `json:"city" example:"São Paulo" binding:"required"`
	Neighborhood string `json:"neighborhood" example:"Centro Histórico" binding:"required"`
	Street       string `json:"street" example:"R. Vinte e Cinco de Março" binding:"required"`
	Number       string `json:"number" example:"1000" binding:"required"`
	Complement   string `json:"complement" example:"1A" binding:"required"`
	Zipcode      string `json:"zipcode" example:"01021100" binding:"required"`
}

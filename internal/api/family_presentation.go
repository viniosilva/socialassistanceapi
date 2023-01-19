package api

type Family struct {
	ID           int    `json:"id" example:"1"`
	Name         string `json:"name" example:"Sauro"`
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

type FamilyResponse struct {
	Data *Family `json:"data"`
}

type FamiliesResponse struct {
	PaginationResponse
	Data []Family `json:"data"`
}

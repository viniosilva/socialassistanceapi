package service

type Resource struct {
	ID          int     `json:"id" example:"1"`
	CreatedAt   string  `json:"created_at" example:"2000-01-01T12:03:00"`
	UpdatedAt   string  `json:"updated_at" example:"2000-01-01T12:03:00"`
	DeletedAt   string  `json:"deleted_at" example:""`
	Name        string  `json:"name" example:"Arroz"`
	Amount      float64 `json:"amount" example:"5"`
	Measurement string  `json:"measurement" example:"Kg"`
	Quantity    float64 `json:"quantity" example:"10"`
}

type ResourceResponse struct {
	Data *Resource `json:"data"`
}

type ResourcesResponse struct {
	Data []Resource `json:"data"`
}

type CreateResourceDto struct {
	Name        string  `json:"name" example:"Arroz" binding:"required"`
	Amount      float64 `json:"amount" example:"5" binding:"required,gte=0"`
	Measurement string  `json:"measurement" example:"Kg" binding:"required"`
	Quantity    float64 `json:"quantity" example:"10" binding:"required,gte=0"`
}

type UpdateResourceDto struct {
	ID          int     `json:"-"`
	Name        string  `json:"name" example:"Arroz"`
	Amount      float64 `json:"amount" example:"5" binding:"gte=0"`
	Measurement string  `json:"measurement" example:"Kg"`
}

type UpdateResourceQuantityDto struct {
	Quantity float64 `json:"quantity" example:"10" binding:"required,gte=0"`
}

package service

type Resource struct {
	ID          int     `json:"id" example:"1"`
	CreatedAt   string  `json:"created_at" example:"2000-01-01T12:03:00"`
	UpdatedAt   string  `json:"updated_at" example:"2000-01-01T12:03:00"`
	DeletedAt   string  `json:"deleted_at" example:""`
	Name        string  `json:"name" example:"Arroz"`
	Amount      float32 `json:"amount" example:"5"`
	Measurement string  `json:"measurement" example:"Kg"`
}

type ResourceResponse struct {
	Data *Resource `json:"data"`
}

type ResourcesResponse struct {
	Data []Resource `json:"data"`
}

type ResourceDto struct {
	Name        string  `json:"name" example:"Arroz" binding:"required"`
	Amount      float32 `json:"amount" example:"5" binding:"required"`
	Measurement string  `json:"measurement" example:"Kg" binding:"required"`
}

type ResourceUpdateDto struct {
	Name        string `json:"name" example:"Arroz"`
	Measurement string `json:"measurement" example:"Kg"`
}

type ResourceTransferAmountDto struct {
	Amount float32 `json:"amount" example:"5" binding:"required"`
}

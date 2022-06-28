package service

type Resource struct {
	ID          int    `json:"id" example:"1"`
	Name        string `json:"name_at" example:"2000-01-01T12:03:00"`
	amount      string `json:"amount_at" example:"2000-01-01T12:03:00"`
	measurement string `json:"measurement_at" example:"2000-01-01T12:03:00"`
	CreatedAt   string `json:"created_at" example:"2000-01-01T12:03:00"`
	UpdatedAt   string `json:"updated_at" example:"2000-01-01T12:03:00"`
	DeletedAt   string `json:"deleted_at" example:"2000-01-01T12:03:00"`
}

type ResourceResponse struct {
	Data *Resource `json:"data"`
}

type ResourcesResponse struct {
	Data []Resource `json:"data"`
}

type ResourceDto struct {
	Name string `json:"name" example:"Cláudio" binding:"required"`
}

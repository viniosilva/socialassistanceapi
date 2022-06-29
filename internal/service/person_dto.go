package service

type Person struct {
	ID        int    `json:"id" example:"1"`
	CreatedAt string `json:"created_at" example:"2000-01-01T12:03:00"`
	UpdatedAt string `json:"updated_at" example:"2000-01-01T12:03:00"`
	DeletedAt string `json:"deleted_at" example:"2000-01-01T12:03:00"`
	Name      string `json:"name" example:"Cláudio"`
}

type PersonResponse struct {
	Data *Person `json:"data"`
}

type PersonsResponse struct {
	Data []Person `json:"data"`
}

type PersonDto struct {
	Name string `json:"name" example:"Cláudio" binding:"required"`
}

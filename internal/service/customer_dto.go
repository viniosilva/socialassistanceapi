package service

type Customer struct {
	ID        int    `json:"id" example:"1"`
	CreatedAt string `json:"created_at" example:"2000-01-01"`
	UpdatedAt string `json:"updated_at" example:"2000-01-01"`
	DeletedAt string `json:"deleted_at" example:"2000-01-01"`
	Name      string `json:"name" example:"Cláudio"`
}

type CustomerResponse struct {
	Data *Customer `json:"data"`
}

type CustomersResponse struct {
	Data []Customer `json:"data"`
}

type CustomerDto struct {
	Name string `json:"name" example:"Cláudio" binding:"required"`
}

package service

type Customer struct {
	ID   int    `json:"id" example:"1"`
	Name string `json:"name" example:"Cláudio"`
}

type CustomerResponse struct {
	Data *Customer `json:"data"`
}

type CustomersResponse struct {
	Data []Customer `json:"data"`
}

type CreateCustomerDto struct {
	Name string `json:"name" example:"Cláudio" binding:"required"`
}

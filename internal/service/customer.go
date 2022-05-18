package service

import (
	"context"

	"github.com/viniosilva/socialassistanceapi/internal/model"
	"github.com/viniosilva/socialassistanceapi/internal/store"
)

type CustomerService struct {
	store store.CustomerStore
}

func NewCustomerService(store store.CustomerStore) *CustomerService {
	return &CustomerService{store}
}

func (impl *CustomerService) FindAll(ctx context.Context) (CustomersResponse, error) {
	customers, err := impl.store.FindAll(ctx)
	if err != nil {
		return CustomersResponse{}, err
	}

	res := []Customer{}
	for _, c := range customers {
		res = append(res, Customer{
			ID:        c.ID,
			CreatedAt: c.CreatedAt.Format("2006-01-02"),
			UpdatedAt: c.UpdatedAt.Format("2006-01-02"),
			Name:      c.Name,
		})
	}

	return CustomersResponse{Data: res}, nil
}

func (impl *CustomerService) FindOneById(ctx context.Context, customerID int) (CustomerResponse, error) {
	customer, err := impl.store.FindOneById(ctx, customerID)
	if err != nil || customer == nil {
		return CustomerResponse{}, err
	}

	return CustomerResponse{
		Data: &Customer{
			ID:        customer.ID,
			CreatedAt: customer.CreatedAt.Format("2006-01-02"),
			UpdatedAt: customer.UpdatedAt.Format("2006-01-02"),
			Name:      customer.Name,
		},
	}, nil
}

func (impl *CustomerService) Create(ctx context.Context, dto CustomerDto) (CustomerResponse, error) {
	customer, err := impl.store.Create(ctx, model.Customer{Name: dto.Name})
	if err != nil {
		return CustomerResponse{}, err
	}

	return CustomerResponse{
		Data: &Customer{
			ID:        customer.ID,
			CreatedAt: customer.CreatedAt.Format("2006-01-02"),
			UpdatedAt: customer.UpdatedAt.Format("2006-01-02"),
			Name:      customer.Name,
		},
	}, nil
}

func (impl *CustomerService) Update(ctx context.Context, customerID int, dto CustomerDto) (CustomerResponse, error) {
	customer, err := impl.store.Update(ctx, model.Customer{
		ID:   customerID,
		Name: dto.Name,
	})
	if err != nil || customer == nil {
		return CustomerResponse{}, err
	}

	return CustomerResponse{
		Data: &Customer{
			ID:        customer.ID,
			CreatedAt: customer.CreatedAt.Format("2006-01-02"),
			UpdatedAt: customer.UpdatedAt.Format("2006-01-02"),
			Name:      customer.Name,
		},
	}, nil
}

func (impl *CustomerService) Delete(ctx context.Context, customerID int) (CustomerResponse, error) {
	err := impl.store.Delete(ctx, customerID)
	return CustomerResponse{}, err
}

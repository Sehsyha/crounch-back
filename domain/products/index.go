package products

import (
	"github.com/crounch-me/back/domain"
	"github.com/crounch-me/back/domain/users"
)

type ProductService struct {
	ProductStorage Storage
	Generation     domain.Generation
}

func (ps *ProductService) GetProduct(productID, userID string) (*Product, *domain.Error) {
	product, err := ps.ProductStorage.GetProduct(productID)
	if err != nil {
		return nil, err
	}

	if !IsUserAuthorized(product, userID) {
		return nil, domain.NewError(domain.UnauthorizedErrorCode)
	}

	return product, nil
}

func (ps *ProductService) CreateProduct(name, userID string) (*Product, *domain.Error) {
	id, err := ps.Generation.GenerateID()
	if err != nil {
		return nil, err
	}

	product := &Product{
		ID:   id,
		Name: name,
		Owner: &users.User{
			ID: userID,
		},
	}

	err = ps.ProductStorage.CreateProduct(product)
	if err != nil {
		return nil, err
	}

	return product, nil
}

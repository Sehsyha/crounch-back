package lists

import (
	"github.com/crounch-me/back/domain"
	"github.com/crounch-me/back/domain/products"
	"github.com/crounch-me/back/domain/users"
)

type ListService struct {
	ListStorage    Storage
	ProductStorage products.Storage
	Generation     domain.Generation
}

func (ls *ListService) CreateList(name, userID string) (*List, *domain.Error) {
	id, err := ls.Generation.GenerateID()
	if err != nil {
		return nil, err
	}

	err = ls.ListStorage.CreateList(id, name, userID)

	if err != nil {
		return nil, err
	}

	list := &List{
		ID:   id,
		Name: name,
		Owner: &users.User{
			ID: userID,
		},
	}

	return list, nil
}

func (ls *ListService) GetOwnersLists(ownerID string) ([]*List, *domain.Error) {
	lists, err := ls.ListStorage.GetOwnersLists(ownerID)

	if err != nil {
		return nil, err
	}

	return lists, nil
}

func (ls *ListService) GetList(listID, userID string) (*List, *domain.Error) {
	list, err := ls.ListStorage.GetList(listID)
	if err != nil {
		return nil, err
	}

	if !IsUserAuthorized(list, userID) {
		return nil, domain.NewError(domain.UnauthorizedErrorCode)
	}

	return list, err
}

func (ls *ListService) AddProductToList(productID, listID, userID string) (*ProductInList, *domain.Error) {
	list, err := ls.ListStorage.GetList(listID)
	if err != nil {
		return nil, err
	}

	if !IsUserAuthorized(list, userID) {
		return nil, domain.NewError(domain.UnauthorizedErrorCode)
	}

	product, err := ls.ProductStorage.GetProduct(productID)
	if err != nil {
		return nil, err
	}

	if !products.IsUserAuthorized(product, userID) {
		return nil, domain.NewError(domain.UnauthorizedErrorCode)
	}

	productInList, err := ls.ListStorage.GetProductInList(productID, listID)
	if err == nil {
		return nil, domain.NewError(DuplicateProductInListErrorCode)
	} else if err.Code != ProductInListNotFoundErrorCode {
		return nil, err
	}

	err = ls.ListStorage.AddProductToList(productID, listID)
	if err != nil {
		return nil, err
	}

	productInList = &ProductInList{
		ProductID: productID,
		ListID:    listID,
	}

	return productInList, nil
}

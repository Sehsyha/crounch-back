package lists

import (
	"time"

	"github.com/crounch-me/back/domain"
)

// ListStorage defines every data functions that we need
type Storage interface {
	CreateList(id, name, ownerID string, creationDate time.Time) *domain.Error
	GetOwnersLists(ownerID string) ([]*List, *domain.Error)
	GetList(id string) (*List, *domain.Error)
	GetProductInList(productID string, listID string) (*ProductInList, *domain.Error)
	AddProductToList(productID string, listID string) *domain.Error
	DeleteList(listID string) *domain.Error
	DeleteProductsFromList(listID string) *domain.Error
}
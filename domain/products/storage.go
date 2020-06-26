package products

import (
	"github.com/crounch-me/back/domain"
)

type Storage interface {
	CreateProduct(id, name, ownerID string) *domain.Error
	GetProduct(id string) (*Product, *domain.Error)
	SearchDefaults(lowerCasedName string, id string) ([]*Product, *domain.Error)
}

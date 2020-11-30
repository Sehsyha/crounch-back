package list

import (
	"time"

	"github.com/crounch-me/back/internal/products"
	"github.com/crounch-me/back/internal/user"
)

// List represents a shopping list
type List struct {
	ID              string           `json:"id"`
	Name            string           `json:"name" validate:"required,lt=61"`
	CreationDate    time.Time        `json:"creationDate"`
	ArchivationDate *time.Time       `json:"archivationDate,omitempty"`
	Contributors    []*user.User     `json:"contributors,omitempty"`
	Products        []*ProductInList `json:"products,omitempty"`
}

// ProductInListLink represents a product in a list
type ProductInListLink struct {
	ProductID string `json:"productId"`
	ListID    string `json:"listId"`
	Bought    bool   `json:"bought"`
}

type ProductInList struct {
	*products.Product
	Bought bool `json:"bought"`
}

// UpdateProductInList represents the possible attributes to update in a product in a list
type UpdateProductInList struct {
	Bought bool `json:"bought"`
}

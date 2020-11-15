package storage

import (
	"github.com/crounch-me/back/domain/authorization.go"
	"github.com/crounch-me/back/domain/contributors"
	"github.com/crounch-me/back/domain/lists"
	"github.com/crounch-me/back/domain/products"
	"github.com/crounch-me/back/domain/users"
)

// Storage defines every data functions that we need
type Storage interface {
	users.Storage
	authorization.Storage
	lists.Storage
  products.Storage
  contributors.Storage
}

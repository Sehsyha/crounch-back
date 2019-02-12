package mock

import (
	"github.com/Sehsyha/crounch-back/model"
	"github.com/Sehsyha/crounch-back/storage"
	testmock "github.com/stretchr/testify/mock"
)

type StorageMock struct {
	testmock.Mock
}

func NewStorageMock() storage.Storage {
	return &StorageMock{}
}

func (sm *StorageMock) CreateUser(u *model.User) error {
	args := sm.Called(u)
	return args.Error(0)
}

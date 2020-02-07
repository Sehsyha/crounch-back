package mock

import "github.com/Sehsyha/crounch-back/model"

// CreateList mocks create list
func (sm *StorageMock) CreateList(list *model.List) error {
	args := sm.Called(list)
	return args.Error(0)
}

// GetOwnerLists mocks get list with owner id
func (sm *StorageMock) GetOwnerLists(ownerID string) ([]*model.List, error) {
	args := sm.Called(ownerID)
	return args.Get(0).([]*model.List), args.Error(1)
}

// GetList mocks get list with an id
func (sm *StorageMock) GetList(id string) (*model.List, error) {
	args := sm.Called(id)
	return args.Get(0).(*model.List), args.Error(1)
}

package users

import (
	"testing"

	"github.com/crounch-me/back/domain"
	"github.com/stretchr/testify/assert"
)

func TestGetByEmailGetByEmailError(t *testing.T) {
	email := "email"

	userStorageMock := &StorageMock{}
	userStorageMock.On("GetByEmail", email).Return(nil, domain.NewError(domain.UnknownErrorCode))

	userService := &UserService{
		UserStorage: userStorageMock,
	}

	result, err := userService.GetByEmail(email)

	userStorageMock.AssertCalled(t, "GetByEmail", email)
	assert.Empty(t, result)
	assert.Equal(t, domain.UnknownErrorCode, err.Code)
}

func TestGetByEmailOK(t *testing.T) {
	email := "email"
	user := &User{
		Email: email,
	}

	userStorageMock := &StorageMock{}
	userStorageMock.On("GetByEmail", email).Return(user, nil)

	userService := &UserService{
		UserStorage: userStorageMock,
	}

	result, err := userService.GetByEmail(email)

	userStorageMock.AssertCalled(t, "GetByEmail", email)
	assert.Equal(t, user, result)
	assert.Empty(t, err)
}
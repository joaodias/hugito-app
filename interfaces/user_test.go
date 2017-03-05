package interfaces_test

import (
	"github.com/joaodias/hugito-backend/domain"
	"github.com/joaodias/hugito-backend/interfaces"
	"github.com/joaodias/hugito-backend/interfaces/mocks"
	"github.com/pborman/uuid"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
	"testing"
)

var user domain.User

func setupUser() {
	user.ID = uuid.New()
	user.Name = "Joao Dias"
	user.Email = "diasjoaoac@gmail.com"
	user.AccessToken = "someToken"
	user.Login = "joaodias"
}

func TestNewUserFail(t *testing.T) {
	setupUser()
	mockDatabaseHandler := &mocks.DatabaseHandler{
		IsError: true,
	}
	mockExternalUserRepository := &interfaces.ExternalUserRepository{
		DatabaseHandler: mockDatabaseHandler,
		Logger:          &mocks.Logger{},
	}
	err := mockExternalUserRepository.New(user)
	assert.True(t, mockDatabaseHandler.AddCalled)
	assert.NotNil(t, err)
}

func TestNewUserSuccess(t *testing.T) {
	setupUser()
	mockDatabaseHandler := &mocks.DatabaseHandler{
		IsError: false,
	}
	mockExternalUserRepository := &interfaces.ExternalUserRepository{
		DatabaseHandler: mockDatabaseHandler,
		Logger:          &mocks.Logger{},
	}
	err := mockExternalUserRepository.New(user)
	assert.True(t, mockDatabaseHandler.AddCalled)
	assert.Nil(t, err)
	assert.Equal(t, mockDatabaseHandler.Table, "user")
	assert.Equal(t, mockDatabaseHandler.Conflict, "update")
}

func TestReadUserFail(t *testing.T) {
	setupUser()
	mockSourceControl := &mocks.SourceControl{
		IsGetUserError: true,
	}
	mockExternalUserRepository := &interfaces.ExternalUserRepository{
		SourceControl: mockSourceControl,
		Logger:        &mocks.Logger{},
	}
	newUser, err := mockExternalUserRepository.Read(user.AccessToken, &oauth2.Config{})
	assert.True(t, mockSourceControl.GetUserCalled)
	assert.NotNil(t, err)
	assert.Nil(t, newUser)
}

func TestReadUserSuccess(t *testing.T) {
	setupUser()
	mockSourceControl := &mocks.SourceControl{}
	mockExternalUserRepository := &interfaces.ExternalUserRepository{
		SourceControl: mockSourceControl,
		Logger:        &mocks.Logger{},
	}
	newUser, err := mockExternalUserRepository.Read(user.AccessToken, &oauth2.Config{})
	assert.True(t, mockSourceControl.GetUserCalled)
	assert.Nil(t, err)
	assert.NotNil(t, newUser)
}

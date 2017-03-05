package usecases_test

import (
	"github.com/joaodias/hugito-app/usecases"
	"github.com/joaodias/hugito-app/usecases/mocks"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
	"testing"
)

type userData struct {
	userName    string
	email       string
	accessToken string
	login       string
}

var user userData

func setupUser() {
	user.userName = "Joao Dias"
	user.email = "diasjoaoac@gmail.com"
	user.accessToken = "someToken"
	user.login = "joaodias"
}

func TestNewUserFail(t *testing.T) {
	setupUser()
	mockUserRepository := &mocks.UserRepository{
		IsError: true,
	}
	userInteractor := &usecases.UserInteractor{
		UserRepository: mockUserRepository,
		Logger:         &mocks.Logger{},
	}
	newUser, err := userInteractor.New(user.userName, user.email, user.login, user.accessToken)
	assert.Nil(t, newUser)
	assert.NotNil(t, err)
	assert.True(t, mockUserRepository.NewCalled)
}

func TestNewUserSuccess(t *testing.T) {
	setupUser()
	mockUserRepository := &mocks.UserRepository{
		IsError: false,
	}
	userInteractor := &usecases.UserInteractor{
		UserRepository: mockUserRepository,
	}
	newUser, err := userInteractor.New(user.userName, user.email, user.login, user.accessToken)
	assert.Nil(t, err, "Error should be nil when the user is successfully created.")
	assert.NotNil(t, newUser)
	assert.True(t, mockUserRepository.NewCalled)
	assert.NotEmpty(t, newUser.ID)
	assert.NotEmpty(t, newUser.AccessToken)
	assert.Equal(t, user.userName, newUser.Name)
	assert.Equal(t, user.login, newUser.Login)
	assert.Equal(t, user.email, newUser.Email)
}

func TestReadUserFail(t *testing.T) {
	setupUser()
	mockUserRepository := &mocks.UserRepository{
		IsError: true,
	}
	userInteractor := &usecases.UserInteractor{
		UserRepository: mockUserRepository,
		Logger:         &mocks.Logger{},
	}
	newUser, err := userInteractor.Read(user.accessToken, &oauth2.Config{})
	assert.True(t, mockUserRepository.ReadCalled)
	assert.Nil(t, newUser)
	assert.NotNil(t, err)
}

func TestReadUserSuccess(t *testing.T) {
	setupUser()
	mockUserRepository := &mocks.UserRepository{
		IsError: false,
	}
	userInteractor := &usecases.UserInteractor{
		UserRepository: mockUserRepository,
		Logger:         &mocks.Logger{},
	}
	newUser, err := userInteractor.Read(user.accessToken, &oauth2.Config{})
	assert.True(t, mockUserRepository.ReadCalled)
	assert.NotNil(t, newUser)
	assert.Nil(t, err)
}

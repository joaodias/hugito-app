package interfaces_test

import (
	"github.com/joaodias/hugito-app/domain"
	"github.com/joaodias/hugito-app/interfaces"
	"github.com/joaodias/hugito-app/interfaces/mocks"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
	"testing"
)

var repository domain.Repository

func setupRepository() {
	repository.ID = "id"
	repository.Name = "name"
	repository.ProjectBranch = "projectBranch"
	repository.PublicBranch = "publicBranch"
	repository.AccessToken = "accessToken"
}

func TestNewRepositoryFail(t *testing.T) {
	setupRepository()
	mockDatabaseHandler := &mocks.DatabaseHandler{
		IsError: true,
	}
	mockExternalRepositoryRepository := &interfaces.ExternalRepositoryRepository{
		DatabaseHandler: mockDatabaseHandler,
		Logger:          &mocks.Logger{},
	}
	err := mockExternalRepositoryRepository.New(repository)
	assert.True(t, mockDatabaseHandler.AddCalled)
	assert.NotNil(t, err)
}

func TestNewRepositorySuccess(t *testing.T) {
	setupRepository()
	mockDatabaseHandler := &mocks.DatabaseHandler{
		IsError: false,
	}
	mockExternalRepositoryRepository := &interfaces.ExternalRepositoryRepository{
		DatabaseHandler: mockDatabaseHandler,
		Logger:          &mocks.Logger{},
	}
	err := mockExternalRepositoryRepository.New(repository)
	assert.True(t, mockDatabaseHandler.AddCalled)
	assert.Nil(t, err)
	assert.Equal(t, "repository", mockDatabaseHandler.Table)
	assert.Equal(t, "update", mockDatabaseHandler.Conflict)
}

func TestValidateRepositoryFail(t *testing.T) {
	setupRepository()
	mockSourceControl := &mocks.SourceControl{
		IsValidateRepositoryError: true,
	}
	mockExternalRepositoryRepository := &interfaces.ExternalRepositoryRepository{
		SourceControl: mockSourceControl,
		Logger:        &mocks.Logger{},
	}
	isValid, err := mockExternalRepositoryRepository.Validate(repository.AccessToken, &oauth2.Config{}, repository)
	assert.True(t, mockSourceControl.ValidateRepositoryCalled)
	assert.NotNil(t, err)
	assert.False(t, isValid)
}

func TestValidateRepositorySuccess(t *testing.T) {
	setupRepository()
	mockSourceControl := &mocks.SourceControl{}
	mockExternalRepositoryRepository := &interfaces.ExternalRepositoryRepository{
		SourceControl: mockSourceControl,
		Logger:        &mocks.Logger{},
	}
	isValid, err := mockExternalRepositoryRepository.Validate(repository.AccessToken, &oauth2.Config{}, repository)
	assert.True(t, mockSourceControl.ValidateRepositoryCalled)
	assert.Nil(t, err)
	assert.True(t, isValid)
}

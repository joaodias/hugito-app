package usecases_test

import (
	"github.com/joaodias/hugito-backend/usecases"
	"github.com/joaodias/hugito-backend/usecases/mocks"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
	"testing"
)

type repositoryData struct {
	id            string
	name          string
	projectBranch string
	publicBranch  string
	accessToken   string
}

var repository repositoryData

func setupRepository() {
	repository.id = "id"
	repository.name = "name"
	repository.projectBranch = "projectBranch"
	repository.publicBranch = "publicBranch"
	repository.accessToken = "accessToken"
}

func TestNewRepositoryFail(t *testing.T) {
	setupRepository()
	mockRepositoryRepository := &mocks.RepositoryRepository{
		IsError: true,
	}
	repositoryInteractor := &usecases.RepositoryInteractor{
		RepositoryRepository: mockRepositoryRepository,
		Logger:               &mocks.Logger{},
	}
	repository, err := repositoryInteractor.New(repository.name, repository.projectBranch, repository.publicBranch, repository.accessToken)
	assert.True(t, mockRepositoryRepository.NewCalled)
	assert.Nil(t, repository)
	assert.NotNil(t, err)
}

func TestNewRepositorySuccess(t *testing.T) {
	setupRepository()
	mockRepositoryRepository := &mocks.RepositoryRepository{
		IsError: false,
	}
	repositoryInteractor := &usecases.RepositoryInteractor{
		RepositoryRepository: mockRepositoryRepository,
		Logger:               &mocks.Logger{},
	}
	newRepository, err := repositoryInteractor.New(repository.name, repository.projectBranch, repository.publicBranch, repository.accessToken)
	assert.True(t, mockRepositoryRepository.NewCalled)
	assert.Nil(t, err)
	assert.NotNil(t, repository)
	assert.Equal(t, repository.name, newRepository.Name)
	assert.Equal(t, repository.projectBranch, newRepository.ProjectBranch)
	assert.Equal(t, repository.publicBranch, newRepository.PublicBranch)
	assert.Equal(t, repository.accessToken, newRepository.AccessToken)
}

func TestValidateRepositoryFail(t *testing.T) {
	setupRepository()
	mockRepositoryRepository := &mocks.RepositoryRepository{
		IsError: true,
	}
	repositoryInteractor := &usecases.RepositoryInteractor{
		RepositoryRepository: mockRepositoryRepository,
		Logger:               &mocks.Logger{},
	}
	isValid, err := repositoryInteractor.Validate(repository.name, repository.projectBranch, repository.accessToken, &oauth2.Config{})
	assert.True(t, mockRepositoryRepository.ValidateCalled)
	assert.False(t, isValid)
	assert.NotNil(t, err)
}

func TestValidateRepositorySuccess(t *testing.T) {
	setupRepository()
	mockRepositoryRepository := &mocks.RepositoryRepository{
		IsError: false,
	}
	repositoryInteractor := &usecases.RepositoryInteractor{
		RepositoryRepository: mockRepositoryRepository,
		Logger:               &mocks.Logger{},
	}
	isValid, err := repositoryInteractor.Validate(repository.name, repository.projectBranch, repository.accessToken, &oauth2.Config{})
	assert.True(t, mockRepositoryRepository.ValidateCalled)
	assert.True(t, isValid)
	assert.Nil(t, err)
}

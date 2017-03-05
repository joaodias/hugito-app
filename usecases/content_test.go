package usecases_test

import (
	"github.com/joaodias/hugito-backend/usecases"
	"github.com/joaodias/hugito-backend/usecases/mocks"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
	"testing"
)

var content usecases.Content

func setupContent() {
	commit := usecases.Commit{
		SHA:   "sha",
		Name:  "name",
		Email: "email",
	}
	content.ID = "id"
	content.RepositoryName = "repository"
	content.ProjectBranch = "projectBranch"
	content.Title = "title"
	content.Path = "path"
	content.Body = "body"
	content.Commit = commit
	content.AccessToken = "accessToken"
}

func TestNewContentFail(t *testing.T) {
	setupContent()
	mockContentRepository := &mocks.ContentRepository{
		IsNewError: true,
	}
	contentInteractor := &usecases.ContentInteractor{
		ContentRepository: mockContentRepository,
	}
	newContent, err := contentInteractor.New(content, &oauth2.Config{})
	assert.True(t, mockContentRepository.NewCalled)
	assert.Nil(t, newContent)
	assert.NotNil(t, err)
}

func TestNewContentSuccess(t *testing.T) {
	setupContent()
	mockContentRepository := &mocks.ContentRepository{}
	contentInteractor := &usecases.ContentInteractor{
		ContentRepository: mockContentRepository,
	}
	newContent, err := contentInteractor.New(content, &oauth2.Config{})
	assert.Nil(t, err)
	assert.NotNil(t, newContent)
	assert.True(t, mockContentRepository.NewCalled)
}

func TestRemoveContentFail(t *testing.T) {
	setupContent()
	mockContentRepository := &mocks.ContentRepository{
		IsRemoveError: true,
	}
	contentInteractor := &usecases.ContentInteractor{
		ContentRepository: mockContentRepository,
	}
	err := contentInteractor.Remove(content, &oauth2.Config{})
	assert.True(t, mockContentRepository.RemoveCalled)
	assert.NotNil(t, err)
}

func TestRemoveContentSuccess(t *testing.T) {
	setupContent()
	mockContentRepository := &mocks.ContentRepository{
		IsRemoveError: false,
	}
	contentInteractor := &usecases.ContentInteractor{
		ContentRepository: mockContentRepository,
	}
	err := contentInteractor.Remove(content, &oauth2.Config{})
	assert.True(t, mockContentRepository.RemoveCalled)
	assert.Nil(t, err)
}

func TestUpdateContenFail(t *testing.T) {
	setupContent()
	mockContentRepository := &mocks.ContentRepository{
		IsUpdateError: true,
	}
	contentInteractor := &usecases.ContentInteractor{
		ContentRepository: mockContentRepository,
	}
	updatedContent, err := contentInteractor.Update(content, &oauth2.Config{})
	assert.True(t, mockContentRepository.UpdateCalled)
	assert.Nil(t, updatedContent)
	assert.NotNil(t, err)
}

func TestUpdateContenSuccess(t *testing.T) {
	setupContent()
	mockContentRepository := &mocks.ContentRepository{
		IsUpdateError: false,
	}
	contentInteractor := &usecases.ContentInteractor{
		ContentRepository: mockContentRepository,
	}
	updatedContent, err := contentInteractor.Update(content, &oauth2.Config{})
	assert.Nil(t, err)
	assert.NotNil(t, updatedContent)
	assert.True(t, mockContentRepository.UpdateCalled)
	assert.NotEmpty(t, updatedContent.ID)
	assert.NotEmpty(t, updatedContent.AccessToken)
	assert.Equal(t, content.RepositoryName, updatedContent.RepositoryName)
	assert.Equal(t, content.ProjectBranch, updatedContent.ProjectBranch)
	assert.Equal(t, content.Title, updatedContent.Title)
	assert.Equal(t, content.Path, updatedContent.Path)
	assert.Equal(t, content.Body, updatedContent.Body)
}

func TestListContentFail(t *testing.T) {
	mockContentRepository := &mocks.ContentRepository{
		IsListError: true,
	}
	contentInteractor := &usecases.ContentInteractor{
		ContentRepository: mockContentRepository,
	}
	contents, err := contentInteractor.List(content, &oauth2.Config{})
	assert.True(t, mockContentRepository.ListCalled)
	assert.NotNil(t, err)
	assert.Nil(t, contents)
}

func TestListContentSuccess(t *testing.T) {
	mockContentRepository := &mocks.ContentRepository{
		IsListError: false,
	}
	contentInteractor := &usecases.ContentInteractor{
		ContentRepository: mockContentRepository,
	}
	contents, err := contentInteractor.List(content, &oauth2.Config{})
	assert.True(t, mockContentRepository.ListCalled)
	assert.Nil(t, err)
	assert.NotNil(t, contents)
}

func TestGetFileContentContentFail(t *testing.T) {
	setupContent()
	mockContentRepository := &mocks.ContentRepository{
		IsFindError: true,
	}
	contentInteractor := &usecases.ContentInteractor{
		ContentRepository: mockContentRepository,
	}
	updatedContent, err := contentInteractor.Find(content, &oauth2.Config{})
	assert.True(t, mockContentRepository.FindCalled)
	assert.Nil(t, updatedContent)
	assert.NotNil(t, err)
}

func TestGetFileContentContentSuccess(t *testing.T) {
	setupContent()
	mockContentRepository := &mocks.ContentRepository{}
	contentInteractor := &usecases.ContentInteractor{
		ContentRepository: mockContentRepository,
	}
	newContent, err := contentInteractor.Find(content, &oauth2.Config{})
	assert.True(t, mockContentRepository.FindCalled)
	assert.Nil(t, err)
	assert.NotNil(t, newContent)
}

func TestPublishContentContentFailWhenUpdatingContent(t *testing.T) {
	setupContent()
	mockContentRepository := &mocks.ContentRepository{
		IsUpdateError: true,
	}
	contentInteractor := &usecases.ContentInteractor{
		ContentRepository: mockContentRepository,
	}
	err := contentInteractor.Publish(content, &oauth2.Config{})
	assert.True(t, mockContentRepository.UpdateCalled)
	assert.False(t, mockContentRepository.PublishCalled)
	assert.NotNil(t, err)
}

func TestPublishContentContentFailWhenPublishingContent(t *testing.T) {
	setupContent()
	mockContentRepository := &mocks.ContentRepository{
		IsPublishError: true,
	}
	contentInteractor := &usecases.ContentInteractor{
		ContentRepository: mockContentRepository,
	}
	err := contentInteractor.Publish(content, &oauth2.Config{})
	assert.True(t, mockContentRepository.UpdateCalled)
	assert.True(t, mockContentRepository.PublishCalled)
	assert.NotNil(t, err)
}

func TestPublishContentContentSuccess(t *testing.T) {
	setupContent()
	mockContentRepository := &mocks.ContentRepository{}
	contentInteractor := &usecases.ContentInteractor{
		ContentRepository: mockContentRepository,
	}
	err := contentInteractor.Publish(content, &oauth2.Config{})
	assert.True(t, mockContentRepository.PublishCalled)
	assert.Nil(t, err)
}

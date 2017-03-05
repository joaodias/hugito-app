package interfaces_test

import (
	"github.com/joaodias/hugito-backend/domain"
	"github.com/joaodias/hugito-backend/interfaces"
	"github.com/joaodias/hugito-backend/interfaces/mocks"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
	"testing"
)

var content domain.Content

func setupContent() {
	commit := domain.Commit{
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

func TestNewContentFailWhenStoringInSourceControl(t *testing.T) {
	setupContent()
	mockSourceControl := &mocks.SourceControl{
		IsCreateContentFileError: true,
	}
	mockExternalContentRepository := &interfaces.ExternalContentRepository{
		SourceControl: mockSourceControl,
		Logger:        &mocks.Logger{},
	}
	content, err := mockExternalContentRepository.New(content, &oauth2.Config{})
	assert.True(t, mockSourceControl.CreateContentFileCalled)
	assert.NotNil(t, err)
	assert.Nil(t, content)
}

func TestNewContentFailWhenStoringInDB(t *testing.T) {
	setupContent()
	mockSourceControl := &mocks.SourceControl{}
	mockDatabaseHandler := &mocks.DatabaseHandler{
		IsError: true,
	}
	mockExternalContentRepository := &interfaces.ExternalContentRepository{
		SourceControl:   mockSourceControl,
		DatabaseHandler: mockDatabaseHandler,
		Logger:          &mocks.Logger{},
	}
	content, err := mockExternalContentRepository.New(content, &oauth2.Config{})
	assert.True(t, mockSourceControl.CreateContentFileCalled)
	assert.True(t, mockDatabaseHandler.AddCalled)
	assert.NotNil(t, err)
	assert.Nil(t, content)
}

func TestNewContentSuccess(t *testing.T) {
	setupContent()
	mockSourceControl := &mocks.SourceControl{}
	mockDatabaseHandler := &mocks.DatabaseHandler{}
	mockExternalContentRepository := &interfaces.ExternalContentRepository{
		SourceControl:   mockSourceControl,
		DatabaseHandler: mockDatabaseHandler,
		Logger:          &mocks.Logger{},
	}
	content, err := mockExternalContentRepository.New(content, &oauth2.Config{})
	assert.True(t, mockSourceControl.CreateContentFileCalled)
	assert.True(t, mockDatabaseHandler.AddCalled)
	assert.Nil(t, err)
	assert.NotNil(t, content)
}

func TestGetFileContentFailWhenStoringInSourceControl(t *testing.T) {
	setupContent()
	mockSourceControl := &mocks.SourceControl{
		IsGetFileContentError: true,
	}
	mockExternalContentRepository := &interfaces.ExternalContentRepository{
		SourceControl: mockSourceControl,
		Logger:        &mocks.Logger{},
	}
	updatedContent, err := mockExternalContentRepository.Find(content, &oauth2.Config{})
	assert.True(t, mockSourceControl.GetFileContentCalled)
	assert.NotNil(t, err)
	assert.Nil(t, updatedContent)
}

func TestGetFileContentSuccess(t *testing.T) {
	setupContent()
	mockSourceControl := &mocks.SourceControl{}
	mockExternalContentRepository := &interfaces.ExternalContentRepository{
		SourceControl: mockSourceControl,
		Logger:        &mocks.Logger{},
	}
	content, err := mockExternalContentRepository.Find(content, &oauth2.Config{})
	assert.True(t, mockSourceControl.GetFileContentCalled)
	assert.Nil(t, err)
	assert.NotNil(t, content)
}

func TestRemoveContentFailWhenStoringInSourceControl(t *testing.T) {
	setupContent()
	mockSourceControl := &mocks.SourceControl{
		IsRemoveContentFileError: true,
	}
	mockExternalContentRepository := &interfaces.ExternalContentRepository{
		SourceControl: mockSourceControl,
		Logger:        &mocks.Logger{},
	}
	err := mockExternalContentRepository.Remove(content, &oauth2.Config{})
	assert.True(t, mockSourceControl.RemoveContentFileCalled)
	assert.NotNil(t, err)
}

func TestRemoveContentFailWhenStoringInDB(t *testing.T) {
	setupContent()
	mockSourceControl := &mocks.SourceControl{}
	mockDatabaseHandler := &mocks.DatabaseHandler{
		IsError: true,
	}
	mockExternalContentRepository := &interfaces.ExternalContentRepository{
		SourceControl:   mockSourceControl,
		DatabaseHandler: mockDatabaseHandler,
		Logger:          &mocks.Logger{},
	}
	err := mockExternalContentRepository.Remove(content, &oauth2.Config{})
	assert.True(t, mockSourceControl.RemoveContentFileCalled)
	assert.True(t, mockDatabaseHandler.RemoveCalled)
	assert.NotNil(t, err)
}

func TestRemoveContentSuccess(t *testing.T) {
	setupContent()
	mockSourceControl := &mocks.SourceControl{}
	mockDatabaseHandler := &mocks.DatabaseHandler{}
	mockExternalContentRepository := &interfaces.ExternalContentRepository{
		SourceControl:   mockSourceControl,
		DatabaseHandler: mockDatabaseHandler,
		Logger:          &mocks.Logger{},
	}
	err := mockExternalContentRepository.Remove(content, &oauth2.Config{})
	assert.True(t, mockSourceControl.RemoveContentFileCalled)
	assert.True(t, mockDatabaseHandler.RemoveCalled)
	assert.Nil(t, err)
}

func TestUpdateContentFailWhenStoringInSourceControl(t *testing.T) {
	setupContent()
	mockSourceControl := &mocks.SourceControl{
		IsUpdateFileContentError: true,
	}
	mockExternalContentRepository := &interfaces.ExternalContentRepository{
		SourceControl: mockSourceControl,
		Logger:        &mocks.Logger{},
	}
	err := mockExternalContentRepository.Update(content, &oauth2.Config{})
	assert.True(t, mockSourceControl.UpdateFileContentCalled)
	assert.NotNil(t, err)
}

func TestUpdateContentSuccess(t *testing.T) {
	setupContent()
	mockSourceControl := &mocks.SourceControl{}
	mockExternalContentRepository := &interfaces.ExternalContentRepository{
		SourceControl: mockSourceControl,
		Logger:        &mocks.Logger{},
	}
	err := mockExternalContentRepository.Update(content, &oauth2.Config{})
	assert.True(t, mockSourceControl.UpdateFileContentCalled)
	assert.Nil(t, err)
}

func TestListFailWhenListingFromSourceControl(t *testing.T) {
	setupContent()
	mockSourceControl := &mocks.SourceControl{
		IsListContentTitlesError: true,
	}
	mockExternalContentRepository := &interfaces.ExternalContentRepository{
		SourceControl: mockSourceControl,
		Logger:        &mocks.Logger{},
	}
	contents, err := mockExternalContentRepository.List(content, &oauth2.Config{})
	assert.True(t, mockSourceControl.ListContentTitlesCalled)
	assert.NotNil(t, err)
	assert.Nil(t, contents)
}

func TestListTitlesSuccess(t *testing.T) {
	setupContent()
	mockSourceControl := &mocks.SourceControl{}
	mockExternalContentRepository := &interfaces.ExternalContentRepository{
		SourceControl: mockSourceControl,
		Logger:        &mocks.Logger{},
	}
	contents, err := mockExternalContentRepository.List(content, &oauth2.Config{})
	assert.True(t, mockSourceControl.ListContentTitlesCalled)
	assert.Nil(t, err)
	assert.NotNil(t, contents)
}

func TestPublishFailWhenDownloadingFiles(t *testing.T) {
	setupContent()
	mockSourceControl := &mocks.SourceControl{
		IsDownloadError: true,
	}
	mockExternalContentRepository := &interfaces.ExternalContentRepository{
		SourceControl: mockSourceControl,
		Logger:        &mocks.Logger{},
	}
	err := mockExternalContentRepository.Publish(content, &oauth2.Config{})
	assert.True(t, mockSourceControl.DownloadContentsCalled)
	assert.NotNil(t, err)
}

func TestPublishFailWhenBuildingSite(t *testing.T) {
	setupContent()
	mockSourceControl := &mocks.SourceControl{}
	mockBuildEngine := &mocks.BuildEngine{
		IsError: true,
	}
	mockExternalContentRepository := &interfaces.ExternalContentRepository{
		SourceControl: mockSourceControl,
		BuildEngine:   mockBuildEngine,
		Logger:        &mocks.Logger{},
	}
	err := mockExternalContentRepository.Publish(content, &oauth2.Config{})
	assert.True(t, mockSourceControl.DownloadContentsCalled)
	assert.True(t, mockBuildEngine.BuildSiteCalled)
	assert.True(t, mockSourceControl.RemoveDownloadedContentsCalled)
	assert.NotNil(t, err)
}

func TestPublishFailWhenPushingFiles(t *testing.T) {
	setupContent()
	mockSourceControl := &mocks.SourceControl{
		IsPushFilesError: true,
	}
	mockBuildEngine := &mocks.BuildEngine{}
	mockExternalContentRepository := &interfaces.ExternalContentRepository{
		SourceControl: mockSourceControl,
		BuildEngine:   mockBuildEngine,
		Logger:        &mocks.Logger{},
	}
	err := mockExternalContentRepository.Publish(content, &oauth2.Config{})
	assert.True(t, mockSourceControl.DownloadContentsCalled)
	assert.True(t, mockBuildEngine.BuildSiteCalled)
	assert.True(t, mockSourceControl.PushFilesCalled)
	assert.True(t, mockSourceControl.RemoveDownloadedContentsCalled)
	assert.NotNil(t, err)
}

func TestPublishFailWhenRemovingDownloadContents(t *testing.T) {
	setupContent()
	mockSourceControl := &mocks.SourceControl{
		IsRemoveDownloadedContentsError: true,
	}
	mockBuildEngine := &mocks.BuildEngine{}
	mockExternalContentRepository := &interfaces.ExternalContentRepository{
		SourceControl: mockSourceControl,
		BuildEngine:   mockBuildEngine,
		Logger:        &mocks.Logger{},
	}
	assert.Panics(t, func() {
		mockExternalContentRepository.Publish(content, &oauth2.Config{})
	})
}

func TestPublishSuccess(t *testing.T) {
	setupContent()
	mockSourceControl := &mocks.SourceControl{}
	mockBuildEngine := &mocks.BuildEngine{}
	mockExternalContentRepository := &interfaces.ExternalContentRepository{
		SourceControl: mockSourceControl,
		BuildEngine:   mockBuildEngine,
		Logger:        &mocks.Logger{},
	}
	err := mockExternalContentRepository.Publish(content, &oauth2.Config{})
	assert.True(t, mockSourceControl.DownloadContentsCalled)
	assert.True(t, mockBuildEngine.BuildSiteCalled)
	assert.True(t, mockSourceControl.PushFilesCalled)
	assert.True(t, mockSourceControl.RemoveDownloadedContentsCalled)
	assert.Nil(t, err)
}

package mocks

import (
	"errors"
	"github.com/joaodias/hugito-backend/domain"
	"golang.org/x/oauth2"
)

type SourceControl struct {
	IsGetUserError                  bool
	IsValidateRepositoryError       bool
	IsListContentTitlesError        bool
	IsCreateContentFileError        bool
	IsRemoveContentFileError        bool
	IsGetFileContentError           bool
	IsUpdateFileContentError        bool
	IsDownloadError                 bool
	IsPushFilesError                bool
	IsRemoveDownloadedContentsError bool
	GetUserCalled                   bool
	ValidateRepositoryCalled        bool
	ListContentTitlesCalled         bool
	CreateContentFileCalled         bool
	RemoveContentFileCalled         bool
	GetFileContentCalled            bool
	UpdateFileContentCalled         bool
	DownloadContentsCalled          bool
	PushFilesCalled                 bool
	RemoveDownloadedContentsCalled  bool
}

func (sc *SourceControl) GetUser(accessToken string, oauthConfiguration *oauth2.Config) (*domain.User, error) {
	sc.GetUserCalled = true
	if sc.IsGetUserError {
		return nil, errors.New("Some error")
	}
	return &domain.User{}, nil
}

func (sc *SourceControl) ValidateRepository(accessToken string, oauthConfiguration *oauth2.Config, repository domain.Repository) (bool, error) {
	sc.ValidateRepositoryCalled = true
	if sc.IsValidateRepositoryError {
		return false, errors.New("Some error")
	}
	return true, nil
}

func (sc *SourceControl) ListContentTitles(content domain.Content, oauthConfiguration *oauth2.Config) ([]string, error) {
	sc.ListContentTitlesCalled = true
	if sc.IsListContentTitlesError {
		return nil, errors.New("Some error")
	}
	return make([]string, 1), nil
}

func (sc *SourceControl) CreateContentFile(content domain.Content, oauthConfiguration *oauth2.Config) (*domain.Content, error) {
	sc.CreateContentFileCalled = true
	if sc.IsCreateContentFileError {
		return nil, errors.New("Some error")
	}
	return &domain.Content{}, nil
}

func (sc *SourceControl) RemoveContentFile(content domain.Content, oauthConfiguration *oauth2.Config) error {
	sc.RemoveContentFileCalled = true
	if sc.IsRemoveContentFileError {
		return errors.New("Some error")
	}
	return nil
}

func (sc *SourceControl) GetFileContent(content domain.Content, oauthConfiguration *oauth2.Config) (*string, error) {
	sc.GetFileContentCalled = true
	if sc.IsGetFileContentError {
		return nil, errors.New("Some error")
	}
	someString := "blabla"
	return &someString, nil
}

func (sc *SourceControl) UpdateFileContent(content domain.Content, oauthConfiguration *oauth2.Config) error {
	sc.UpdateFileContentCalled = true
	if sc.IsUpdateFileContentError {
		return errors.New("Some error")
	}
	return nil
}

func (sc *SourceControl) DownloadContents(content domain.Content, oauthConfiguration *oauth2.Config) (*string, error) {
	sc.DownloadContentsCalled = true
	if sc.IsDownloadError {
		return nil, errors.New("Some error")
	}
	someString := "blabla"
	return &someString, nil
}

func (sc *SourceControl) PushFiles(content domain.Content, oauthConfiguration *oauth2.Config, sourcePath string) error {
	sc.PushFilesCalled = true
	if sc.IsPushFilesError {
		return errors.New("Some error")
	}
	return nil
}

func (sc *SourceControl) RemoveDownloadedContents(sourcePath string) error {
	sc.RemoveDownloadedContentsCalled = true
	if sc.IsRemoveDownloadedContentsError {
		return errors.New("Some error")
	}
	return nil
}

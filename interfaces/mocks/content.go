package mocks

import (
	"errors"
	"github.com/joaodias/hugito-app/usecases"
	"golang.org/x/oauth2"
)

type ContentInteractor struct {
	IsNewError     bool
	IsUpdateError  bool
	IsRemoveError  bool
	IsListError    bool
	IsFindError    bool
	IsPublishError bool
	NewCalled      bool
	UpdateCalled   bool
	RemoveCalled   bool
	ListCalled     bool
	FindCalled     bool
	PublishCalled  bool
}

func (ci *ContentInteractor) New(content usecases.Content, oauthConfiguration *oauth2.Config) (*usecases.Content, error) {
	ci.NewCalled = true
	if ci.IsNewError {
		return nil, errors.New("Some error")
	}
	return &usecases.Content{}, nil
}

func (ci *ContentInteractor) Remove(content usecases.Content, oauthConfiguration *oauth2.Config) error {
	ci.RemoveCalled = true
	if ci.IsRemoveError {
		return errors.New("Some error")
	}
	return nil
}

func (ci *ContentInteractor) Update(content usecases.Content, oauthConfiguration *oauth2.Config) (*usecases.Content, error) {
	ci.UpdateCalled = true
	if ci.IsUpdateError {
		return nil, errors.New("Some error")
	}
	return &usecases.Content{}, nil
}

func (ci *ContentInteractor) List(content usecases.Content, oauthConfiguration *oauth2.Config) ([]usecases.Content, error) {
	ci.ListCalled = true
	if ci.IsListError {
		return nil, errors.New("Some error")
	}
	return make([]usecases.Content, 1), nil
}

func (ci *ContentInteractor) Find(content usecases.Content, oauthConfiguration *oauth2.Config) (*usecases.Content, error) {
	ci.FindCalled = true
	if ci.IsFindError {
		return nil, errors.New("Some error")
	}
	return &usecases.Content{}, nil
}

func (ci *ContentInteractor) Publish(content usecases.Content, oauthConfiguration *oauth2.Config) error {
	ci.PublishCalled = true
	if ci.IsPublishError {
		return errors.New("Some error")
	}
	return nil
}

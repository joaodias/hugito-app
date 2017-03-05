package mocks

import (
	"errors"
	"github.com/joaodias/hugito-backend/domain"
	"golang.org/x/oauth2"
)

type ContentRepository struct {
	IsNewError     bool
	IsRemoveError  bool
	IsUpdateError  bool
	IsListError    bool
	IsFindError    bool
	IsPublishError bool
	NewCalled      bool
	RemoveCalled   bool
	UpdateCalled   bool
	ListCalled     bool
	FindCalled     bool
	PublishCalled  bool
}

func (cr *ContentRepository) New(content domain.Content, oauthConfiguration *oauth2.Config) (*domain.Content, error) {
	cr.NewCalled = true
	if cr.IsNewError {
		return nil, errors.New("Some Error")
	}
	return &domain.Content{}, nil
}

func (cr *ContentRepository) Remove(content domain.Content, oauthConfiguration *oauth2.Config) error {
	cr.RemoveCalled = true
	if cr.IsRemoveError {
		return errors.New("Some Error")
	}
	return nil
}

func (cr *ContentRepository) Update(content domain.Content, oauthConfiguration *oauth2.Config) error {
	cr.UpdateCalled = true
	if cr.IsUpdateError {
		return errors.New("Some error")
	}
	return nil
}

func (cr *ContentRepository) List(content domain.Content, oauthConfiguration *oauth2.Config) ([]domain.Content, error) {
	cr.ListCalled = true
	if cr.IsListError {
		return nil, errors.New("Some error")
	}
	return nil, nil
}

func (cr *ContentRepository) Find(content domain.Content, oauthConfiguration *oauth2.Config) (*domain.Content, error) {
	cr.FindCalled = true
	if cr.IsFindError {
		return nil, errors.New("Some error")
	}
	return &domain.Content{}, nil
}

func (cr *ContentRepository) Publish(content domain.Content, oauthConfiguration *oauth2.Config) error {
	cr.PublishCalled = true
	if cr.IsPublishError {
		return errors.New("Some error")
	}
	return nil
}

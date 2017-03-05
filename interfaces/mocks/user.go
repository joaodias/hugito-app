package mocks

import (
	"errors"
	"github.com/joaodias/hugito-backend/usecases"
	"golang.org/x/oauth2"
)

type UserInteractor struct {
	IsReadError  bool
	IsNewError   bool
	IsNewCalled  bool
	IsReadCalled bool
}

func (ui *UserInteractor) New(name, email, login, accessToken string) (*usecases.User, error) {
	ui.IsNewCalled = true
	if ui.IsNewError {
		return nil, errors.New("Some error")
	}
	return &usecases.User{}, nil
}

func (ui *UserInteractor) Read(accessToken string, oauthConfiguration *oauth2.Config) (*usecases.User, error) {
	ui.IsReadCalled = true
	if ui.IsReadError {
		return nil, errors.New("Some error")
	}
	return &usecases.User{}, nil
}

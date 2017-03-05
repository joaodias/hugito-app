package mocks

import (
	"errors"
	"github.com/joaodias/hugito-app/domain"
	"golang.org/x/oauth2"
)

type UserRepository struct {
	NewCalled  bool
	ReadCalled bool
	IsError    bool
}

func (ur *UserRepository) New(user domain.User) error {
	ur.NewCalled = true
	if ur.IsError {
		return errors.New("Some Error")
	}
	return nil
}

func (ur *UserRepository) Read(accessToken string, oauthConfiguration *oauth2.Config) (*domain.User, error) {
	ur.ReadCalled = true
	if ur.IsError {
		return nil, errors.New("Some error")
	}
	return &domain.User{}, nil
}

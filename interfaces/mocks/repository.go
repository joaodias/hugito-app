package mocks

import (
	"errors"
	"github.com/joaodias/hugito-app/usecases"
	"golang.org/x/oauth2"
)

type RepositoryInteractor struct {
	IsNewError        bool
	IsValidateError   bool
	IsNewCalled       bool
	IsValidateCalled  bool
	IsValidRepository bool
}

func (ri *RepositoryInteractor) New(name, projectBranch, publicBranch, accessToken string) (*usecases.Repository, error) {
	ri.IsNewCalled = true
	if ri.IsNewError {
		return nil, errors.New("Some error")
	}
	return &usecases.Repository{}, nil
}

func (ri *RepositoryInteractor) Validate(name string, projectBranch string, accessToken string, oauthConfiguration *oauth2.Config) (bool, error) {
	ri.IsValidateCalled = true
	if ri.IsValidateError {
		return false, errors.New("Some error")
	}
	return ri.IsValidRepository, nil
}

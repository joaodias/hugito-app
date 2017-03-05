package mocks

import (
	"errors"
	"github.com/joaodias/hugito-backend/domain"
	"golang.org/x/oauth2"
)

type RepositoryRepository struct {
	NewCalled      bool
	ValidateCalled bool
	IsError        bool
}

func (rr *RepositoryRepository) New(repository domain.Repository) error {
	rr.NewCalled = true
	if rr.IsError {
		return errors.New("Some error")
	}
	return nil
}

func (rr *RepositoryRepository) Validate(accessToken string, oauthConfiguration *oauth2.Config, repository domain.Repository) (bool, error) {
	rr.ValidateCalled = true
	if rr.IsError {
		return false, errors.New("Some error")
	}
	return true, nil
}

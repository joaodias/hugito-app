package interfaces

import (
	"github.com/joaodias/hugito-backend/domain"
	"golang.org/x/oauth2"
)

type ExternalRepositoryRepository ExternalRepository

func (exr *ExternalRepositoryRepository) New(repository domain.Repository) error {
	err := exr.DatabaseHandler.Add(repository, "repository", "update")
	if err != nil {
		exr.Logger.Log(err.Error())
		return err
	}
	return nil
}

func (exr *ExternalRepositoryRepository) Validate(accessToken string, oauthConfiguration *oauth2.Config, repository domain.Repository) (bool, error) {
	isValid, err := exr.SourceControl.ValidateRepository(accessToken, oauthConfiguration, repository)
	if err != nil {
		exr.Logger.Log(err.Error())
		return false, err
	}
	return isValid, nil
}

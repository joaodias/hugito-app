package interfaces

import (
	"github.com/joaodias/hugito-app/domain"
	"golang.org/x/oauth2"
)

// ExternalUserRepository is the implementation of a
type ExternalUserRepository ExternalRepository

// New is the implementation of New for the user entity. It creates a new user.
func (eur *ExternalUserRepository) New(user domain.User) error {
	err := eur.DatabaseHandler.Add(user, "user", "update")
	if err != nil {
		eur.Logger.Log(err.Error())
		return err
	}
	return nil
}

// Read is the implementation of Read for the user entity. It reads a user from an external source.
func (eur *ExternalUserRepository) Read(accessToken string, oauthConfiguration *oauth2.Config) (*domain.User, error) {
	user, err := eur.SourceControl.GetUser(accessToken, oauthConfiguration)
	if err != nil {
		eur.Logger.Log(err.Error())
		return nil, err
	}
	return user, nil
}

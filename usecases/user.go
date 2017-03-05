package usecases

import (
	"github.com/joaodias/hugito-backend/domain"
	"github.com/pborman/uuid"
	"golang.org/x/oauth2"
	"time"
)

// User is the user of the application
type User struct {
	ID          string    `json:"id" gorethink:"id"`
	CreatedAt   time.Time `json:"createdAt" gorethink:"createdAt"`
	Name        string    `json:"name" gorethink:"name"`
	Email       string    `json:"email" gorethink:"email"`
	Login       string    `json:"login" gorethink:"login"`
	AccessToken string    `json:"accessToken" gorethink:"accessToken"`
}

// UserInteractor is responsible for messing up with the user domain
type UserInteractor struct {
	UserRepository domain.UserRepository
	Logger         Logger
}

// New creates a new user
func (ui *UserInteractor) New(name string, email string, login string, accessToken string) (*User, error) {
	domainUser := domain.User{
		ID:          uuid.New(),
		Name:        name,
		Email:       email,
		Login:       login,
		AccessToken: accessToken,
	}
	err := ui.UserRepository.New(domainUser)
	if err != nil {
		ui.Logger.Log(err.Error())
		return nil, err
	}
	return &User{
		ID:          domainUser.ID,
		Name:        domainUser.Name,
		Email:       domainUser.Email,
		Login:       domainUser.Login,
		AccessToken: domainUser.AccessToken,
	}, nil
}

// Read reads an user
func (ui *UserInteractor) Read(accessToken string, oauthConfiguration *oauth2.Config) (*User, error) {
	domainUser, err := ui.UserRepository.Read(accessToken, oauthConfiguration)
	if err != nil {
		ui.Logger.Log(err.Error())
		return nil, err
	}
	return &User{
		ID:          domainUser.ID,
		Name:        domainUser.Name,
		Email:       domainUser.Email,
		Login:       domainUser.Login,
		AccessToken: domainUser.AccessToken,
	}, nil
}

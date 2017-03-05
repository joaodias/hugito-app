package domain

import (
	"golang.org/x/oauth2"
	"time"
)

// UserRepository abstrats the storage for the user.
type UserRepository interface {
	New(User) error
	Read(accessToken string, oauthConfiguration *oauth2.Config) (*User, error)
}

// User represents the user information exchanged between the server and the client.
type User struct {
	ID          string    `json:"id" gorethink:"id"`
	CreatedAt   time.Time `json:"createdAt" gorethink:"createdAt"`
	Name        string    `json:"name" gorethink:"name"`
	Email       string    `json:"email" gorethink:"email"`
	Login       string    `json:"login" gorethink:"login"`
	AccessToken string    `json:"accessToken" gorethink:"accessToken"`
}

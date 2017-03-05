package domain

import (
	"golang.org/x/oauth2"
	"time"
)

// RepositoryRepository abstracts the storage of repository entries.
type RepositoryRepository interface {
	New(Repository) error
	Validate(string, *oauth2.Config, Repository) (bool, error)
}

// Repository refers to the repository in which the website is stored.
type Repository struct {
	ID            string    `json:"id" gorethink:"id"`
	CreatedAt     time.Time `json:"createdAt" gorethink:"createdAt"`
	Name          string    `json:"name" gorethink:"name"`
	ProjectBranch string    `json:"projectBranch" gorethink:"projectBranch"`
	PublicBranch  string    `json:"publicBranch" gorethink:"publicBranch"`
	AccessToken   string    `json:"accessToken" gorethink:"accessToken"`
}

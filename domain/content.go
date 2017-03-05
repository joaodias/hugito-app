package domain

import (
	"golang.org/x/oauth2"
	"time"
)

// ContentRepository abstracts the storage of content entries.
type ContentRepository interface {
	New(content Content, oauthConfiguration *oauth2.Config) (*Content, error)
	Remove(content Content, oauthConfiguration *oauth2.Config) error
	Update(content Content, oauthConfiguration *oauth2.Config) error
	List(content Content, oauthConfiguration *oauth2.Config) ([]Content, error)
	Find(content Content, oauthConfiguration *oauth2.Config) (*Content, error)
	Publish(content Content, oauthConfiguration *oauth2.Config) error
}

// Content refers to a content entry in the website. It is in a given repository and in a given project branch.
type Content struct {
	ID             string    `json:"id" gorethink:"id"`
	CreatedAt      time.Time `json:"createdAt" gorethink:"createdAt"`
	RepositoryName string    `json:"repositoryName" gorethink:"repositoryName"`
	ProjectBranch  string    `json:"projectBranch" gorethink:"projectBranch"`
	PublicBranch   string    `json:"publicBranch" gorethink:"publicBranch"`
	Title          string    `json:"title" gorethink:"title"`
	Path           string    `json:"path" gorethink:"path"`
	Body           string    `json:"body" gorethink:"body"`
	Commit         Commit    `json:"commit" gorethink:"commit"`
	AccessToken    string    `json:"accessToken" gorethink:"accessToken"`
}

// Commit is a point of commitment in a source control platform.
type Commit struct {
	SHA   string `json:"sha" gorethink:"sha"`
	Name  string `json:"name" gorethink:"name"`
	Email string `json:"email" gorethink:"email"`
}

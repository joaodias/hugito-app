package usecases

import (
	"github.com/joaodias/hugito-backend/domain"
	"github.com/pborman/uuid"
	"golang.org/x/oauth2"
	"time"
)

// Content represents a content in a given repository's branch
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

// ContentInteractor messes with the domain content entity
type ContentInteractor struct {
	ContentRepository domain.ContentRepository
}

func (ci *ContentInteractor) New(content Content, oautConfiguration *oauth2.Config) (*Content, error) {
	domainCommit := domain.Commit{
		SHA:   content.Commit.SHA,
		Name:  content.Commit.Name,
		Email: content.Commit.Email,
	}
	domainContent := domain.Content{
		ID:             uuid.New(),
		RepositoryName: content.RepositoryName,
		ProjectBranch:  content.ProjectBranch,
		PublicBranch:   content.PublicBranch,
		Title:          content.Title,
		Path:           content.Path,
		Body:           content.Body,
		Commit:         domainCommit,
		AccessToken:    content.AccessToken,
	}
	createdContent, err := ci.ContentRepository.New(domainContent, oautConfiguration)
	if err != nil {
		return nil, err
	}
	newCommit := Commit{
		SHA:   createdContent.Commit.SHA,
		Name:  createdContent.Commit.Name,
		Email: createdContent.Commit.Email,
	}
	return &Content{
		ID:             createdContent.ID,
		AccessToken:    createdContent.AccessToken,
		RepositoryName: createdContent.RepositoryName,
		ProjectBranch:  createdContent.ProjectBranch,
		PublicBranch:   createdContent.PublicBranch,
		Title:          createdContent.Title,
		Path:           createdContent.Path,
		Commit:         newCommit,
		Body:           createdContent.Body,
	}, nil
}

func (ci *ContentInteractor) Remove(content Content, oauthConfiguration *oauth2.Config) error {
	domainCommit := domain.Commit{
		SHA:   content.Commit.SHA,
		Name:  content.Commit.Name,
		Email: content.Commit.Email,
	}
	domainContent := domain.Content{
		ID:             content.ID,
		RepositoryName: content.RepositoryName,
		ProjectBranch:  content.ProjectBranch,
		PublicBranch:   content.PublicBranch,
		Title:          content.Title,
		Path:           content.Path,
		Body:           content.Body,
		Commit:         domainCommit,
		AccessToken:    content.AccessToken,
	}
	err := ci.ContentRepository.Remove(domainContent, oauthConfiguration)
	if err != nil {
		return err
	}
	return nil
}

func (ci *ContentInteractor) Update(content Content, oauthConfiguration *oauth2.Config) (*Content, error) {
	domainContent := domain.Content{
		ID:             content.ID,
		RepositoryName: content.RepositoryName,
		ProjectBranch:  content.ProjectBranch,
		PublicBranch:   content.PublicBranch,
		Title:          content.Title,
		Path:           content.Path,
		Body:           content.Body,
		AccessToken:    content.AccessToken,
	}
	err := ci.ContentRepository.Update(domainContent, oauthConfiguration)
	if err != nil {
		return nil, err
	}
	return &Content{
		ID:             domainContent.ID,
		AccessToken:    domainContent.AccessToken,
		RepositoryName: domainContent.RepositoryName,
		ProjectBranch:  domainContent.ProjectBranch,
		PublicBranch:   domainContent.PublicBranch,
		Title:          domainContent.Title,
		Path:           domainContent.Path,
		Body:           domainContent.Body,
	}, nil
}

func (ci *ContentInteractor) List(content Content, oauthConfiguration *oauth2.Config) ([]Content, error) {
	domainCommit := domain.Commit{
		SHA:   content.Commit.SHA,
		Name:  content.Commit.Name,
		Email: content.Commit.Email,
	}
	domainContent := domain.Content{
		ID:             content.ID,
		RepositoryName: content.RepositoryName,
		ProjectBranch:  content.ProjectBranch,
		PublicBranch:   content.PublicBranch,
		Title:          content.Title,
		Path:           content.Path,
		Body:           content.Body,
		Commit:         domainCommit,
		AccessToken:    content.AccessToken,
	}
	domainContents, err := ci.ContentRepository.List(domainContent, oauthConfiguration)
	if err != nil {
		return nil, err
	}
	contents := make([]Content, len(domainContents))
	for i, domainContent := range domainContents {
		contents[i] = Content{
			ID:             domainContent.ID,
			AccessToken:    domainContent.AccessToken,
			RepositoryName: domainContent.RepositoryName,
			ProjectBranch:  domainContent.ProjectBranch,
			PublicBranch:   domainContent.PublicBranch,
			Title:          domainContent.Title,
			Path:           domainContent.Path,
			Body:           domainContent.Body,
		}
	}
	return contents, nil
}

func (ci *ContentInteractor) Find(content Content, oauthConfiguration *oauth2.Config) (*Content, error) {
	domainCommit := domain.Commit{
		SHA:   content.Commit.SHA,
		Name:  content.Commit.Name,
		Email: content.Commit.Email,
	}
	domainContent := domain.Content{
		ID:             content.ID,
		RepositoryName: content.RepositoryName,
		ProjectBranch:  content.ProjectBranch,
		PublicBranch:   content.PublicBranch,
		Title:          content.Title,
		Path:           content.Path,
		Body:           content.Body,
		Commit:         domainCommit,
		AccessToken:    content.AccessToken,
	}
	createdContent, err := ci.ContentRepository.Find(domainContent, oauthConfiguration)
	if err != nil {
		return nil, err
	}
	newCommit := Commit{
		SHA:   createdContent.Commit.SHA,
		Name:  createdContent.Commit.Name,
		Email: createdContent.Commit.Email,
	}
	return &Content{
		ID:             createdContent.ID,
		AccessToken:    createdContent.AccessToken,
		RepositoryName: createdContent.RepositoryName,
		ProjectBranch:  createdContent.ProjectBranch,
		PublicBranch:   createdContent.PublicBranch,
		Title:          createdContent.Title,
		Path:           createdContent.Path,
		Commit:         newCommit,
		Body:           createdContent.Body,
	}, nil
}

func (ci *ContentInteractor) Publish(content Content, oauthConfiguration *oauth2.Config) error {
	domainCommit := domain.Commit{
		SHA:   content.Commit.SHA,
		Name:  content.Commit.Name,
		Email: content.Commit.Email,
	}
	domainContent := domain.Content{
		ID:             content.ID,
		RepositoryName: content.RepositoryName,
		ProjectBranch:  content.ProjectBranch,
		PublicBranch:   content.PublicBranch,
		Title:          content.Title,
		Path:           content.Path,
		Body:           content.Body,
		Commit:         domainCommit,
		AccessToken:    content.AccessToken,
	}
	err := ci.ContentRepository.Update(domainContent, oauthConfiguration)
	if err != nil {
		return err
	}
	err = ci.ContentRepository.Publish(domainContent, oauthConfiguration)
	if err != nil {
		return err
	}
	return nil
}

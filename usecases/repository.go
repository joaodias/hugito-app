package usecases

import (
	"github.com/joaodias/hugito-app/domain"
	"github.com/pborman/uuid"
	"golang.org/x/oauth2"
	"time"
)

// Repository refers to the repository in which the website is stored.
type Repository struct {
	ID            string    `json:"id" gorethink:"id"`
	CreatedAt     time.Time `json:"createdAt" gorethink:"createdAt"`
	Name          string    `json:"name" gorethink:"name"`
	ProjectBranch string    `json:"projectBranch" gorethink:"projectBranch"`
	PublicBranch  string    `json:"publicBranch" gorethink:"publicBranch"`
	AccessToken   string    `json:"accessToken" gorethink:"accessToken"`
}

// RepositoryInteractor messes with the respository domain entity
type RepositoryInteractor struct {
	RepositoryRepository domain.RepositoryRepository
	Logger               Logger
}

// New creates a new repository
func (ri *RepositoryInteractor) New(name, projectBranch, publicBranch, accessToken string) (*Repository, error) {
	domainRepository := domain.Repository{
		ID:            uuid.New(),
		Name:          name,
		ProjectBranch: projectBranch,
		PublicBranch:  publicBranch,
		AccessToken:   accessToken,
	}
	err := ri.RepositoryRepository.New(domainRepository)
	if err != nil {
		return nil, err
	}
	return &Repository{
		ID:            domainRepository.ID,
		Name:          domainRepository.Name,
		ProjectBranch: domainRepository.ProjectBranch,
		PublicBranch:  domainRepository.PublicBranch,
		AccessToken:   domainRepository.AccessToken,
	}, nil
}

// Validate checks the validity of a given repository. The repository must be a valid project. Otherwise it can't be
// used.
func (ri *RepositoryInteractor) Validate(name string, projectBranch string, accessToken string, oauthConfiguration *oauth2.Config) (bool, error) {
	domainRepository := domain.Repository{
		Name:          name,
		ProjectBranch: projectBranch,
		AccessToken:   accessToken,
	}
	isValid, err := ri.RepositoryRepository.Validate(accessToken, oauthConfiguration, domainRepository)
	if err != nil {
		return false, err
	}
	return isValid, nil
}

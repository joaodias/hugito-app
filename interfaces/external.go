package interfaces

import (
	"github.com/joaodias/hugito-backend/domain"
	"golang.org/x/oauth2"
)

// DatabaseHandler is a database abstraction. The interface can get broader or narrower on demand.
type DatabaseHandler interface {
	Add(data interface{}, table string, onConflict string) error
	Remove(id string, table string) error
	Update(data interface{}, id string, table string) error
	List(key string, table string) ([]domain.Content, error)
}

// SourceControl is a source control platform abstraction. It can be github, bitbucket...
type SourceControl interface {
	GetUser(accessToken string, oauthConfiguration *oauth2.Config) (*domain.User, error)
	ValidateRepository(accessToken string, oauthConfiguration *oauth2.Config, repository domain.Repository) (bool, error)
	ListContentTitles(content domain.Content, oauthConfiguration *oauth2.Config) ([]string, error)
	CreateContentFile(content domain.Content, oauthConfiguration *oauth2.Config) (*domain.Content, error)
	RemoveContentFile(content domain.Content, oauthConfiguration *oauth2.Config) error
	GetFileContent(content domain.Content, oauthConfiguration *oauth2.Config) (*string, error)
	UpdateFileContent(content domain.Content, oauthConfiguration *oauth2.Config) error
	DownloadContents(content domain.Content, oauthConfiguration *oauth2.Config) (*string, error)
	PushFiles(content domain.Content, oauthConfiguration *oauth2.Config, sourcePath string) error
	RemoveDownloadedContents(sourcePath string) error
}

// BuildEngine is the static website generator build engine. The engine will build the project branch contents and
// create the public branch contents.
type BuildEngine interface {
	BuildSite(source string) error
}

// ExternalRepository is an abstraction for an external storage. It includes all the necessary interfaces with external
// platforms/services.
type ExternalRepository struct {
	DatabaseHandler DatabaseHandler
	SourceControl   SourceControl
	BuildEngine     BuildEngine
	Logger          Logger
}

package interfaces

import (
	"github.com/joaodias/hugito-backend/domain"
	"golang.org/x/oauth2"
)

// ExternalContentRepository has the logic fro interfacing external content repositories with the application logic.
type ExternalContentRepository ExternalRepository

// New creates a new content file. It creates the file both in the source control and in the database.
func (ecr *ExternalContentRepository) New(content domain.Content, oauthConfiguration *oauth2.Config) (*domain.Content, error) {
	createdContent, err := ecr.SourceControl.CreateContentFile(content, oauthConfiguration)
	if err != nil {
		return nil, err
	}
	err = ecr.DatabaseHandler.Add(content, "content", "update")
	if err != nil {
		return nil, err
	}
	return createdContent, nil
}

// Remove removes a content file both in the source control and in the database.
func (ecr *ExternalContentRepository) Remove(content domain.Content, oauthConfiguration *oauth2.Config) error {
	err := ecr.SourceControl.RemoveContentFile(content, oauthConfiguration)
	if err != nil {
		return err
	}
	err = ecr.DatabaseHandler.Remove(content.ID, "content")
	if err != nil {
		return err
	}
	return nil
}

// Update updates a content file both in the source control and in the database.
func (ecr *ExternalContentRepository) Update(content domain.Content, oauthConfiguration *oauth2.Config) error {
	err := ecr.SourceControl.UpdateFileContent(content, oauthConfiguration)
	if err != nil {
		return err
	}
	return nil
}

// List lists the source control content for a given repository, branch and path. The content fetched from the source
// control is stored in a database in order to get application specific ids to exchange with the client.
func (ecr *ExternalContentRepository) List(content domain.Content, oauthConfiguration *oauth2.Config) ([]domain.Content, error) {
	titles, err := ecr.SourceControl.ListContentTitles(content, oauthConfiguration)
	if err != nil {
		return nil, err
	}
	// Everytime there is a list operation, the db is synced with the source control. This can be improved in the
	// future.
	contents := make([]domain.Content, len(titles))
	for i, title := range titles {
		content.Title = title
		contents[i] = content
	}
	return contents, nil
}

// Find gets the content of a file.
func (ecr *ExternalContentRepository) Find(content domain.Content, oauthConfiguration *oauth2.Config) (*domain.Content, error) {
	fileContent, err := ecr.SourceControl.GetFileContent(content, oauthConfiguration)
	if err != nil {
		return nil, err
	}
	content.Body = *fileContent
	return &content, nil
}

// Publish publishes the content to the source control.
func (ecr *ExternalContentRepository) Publish(content domain.Content, oauthConfiguration *oauth2.Config) error {
	contentsPath, err := ecr.SourceControl.DownloadContents(content, oauthConfiguration)
	if err != nil {
		return err
	}
	err = ecr.BuildEngine.BuildSite(*contentsPath)
	if err != nil {
		removeErr := ecr.SourceControl.RemoveDownloadedContents(*contentsPath)
		if removeErr != nil {
			panic("Can't Remove downloaded contents!")
		}
		return err
	}
	err = ecr.SourceControl.PushFiles(content, oauthConfiguration, *contentsPath+"/public")
	if err != nil {
		removeErr := ecr.SourceControl.RemoveDownloadedContents(*contentsPath)
		if removeErr != nil {
			panic("Can't Remove downloaded contents!")
		}
		return err
	}
	err = ecr.SourceControl.RemoveDownloadedContents(*contentsPath)
	if err != nil {
		panic("Can't Remove downloaded contents!")
	}
	return nil
}

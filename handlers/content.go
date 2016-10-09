package handlers

import (
	"github.com/mitchellh/mapstructure"
)

// ContentList represents the list of content for a given repository name
type ContentList struct {
	Name        string   `json:"name"`
	Titles      []string `json:"title"`
	AccessToken string   `json:"accessToken"`
}

// TODO: Update this structure to take more values and support what I want to the content file options. Things like branch and commit message are important in this structure.

// Content represents the information exchanged between the server and the client.
type Content struct {
	RepositoryName string `json:"repositoryName"`
	Title          string `json:"title"`
	Body           string `json:"body"`
	CommitMessage  string `json:"commitMessage"`
	Branch         string `json:"branch"`
	AccessToken    string `json:"accessToken"`
}

// GetContentList gets the repository content list.
//
// Happy Path:
// 1. Decode JSON
// 2. Get the github authenticated user
// 3. Get the repository content files
// 4. Send the content list to the client
func GetContentList(communicator Communicator, data interface{}) {
	var contentList ContentList
	err := mapstructure.Decode(data, &contentList)
	if err != nil {
		communicator.SetSend("error", "Error decoding json:"+err.Error())
		return
	}
	communicator.NewFinishedChannel(ContentFinished)
	go func() {
		githubClient := GetGithubClient(contentList.AccessToken, communicator)
		userLogin, err := GetGithubUserLogin(githubClient)
		if err != nil {
			communicator.SetSend("logout", "Can't retrieve the authenticated user.")
			communicator.Finished(ContentFinished)
			return
		}
		contentList.Titles, err = GetGithubRepositoryTree(githubClient, userLogin, contentList.Name, "content")
		if err != nil {
			communicator.SetSend("error", "Can't retrieve the content list.")
			communicator.Finished(ContentFinished)
			return
		}
		communicator.SetSend("content list", contentList)
		communicator.Finished(ContentFinished)
	}()
}

// GetFileContent gets the content of a file from the github repository.
//
// Happy Path:
// 1. Decode JSON
// 2. Get the github authenticated user
// 3. Get the content of the given repository file
// 4. Send the content to the client
func GetFileContent(communicator Communicator, data interface{}) {
	var content Content
	err := mapstructure.Decode(data, &content)
	if err != nil {
		communicator.SetSend("error", "Error decoding json:"+err.Error())
		return
	}
	communicator.NewFinishedChannel(FileContentFinished)
	go func() {
		githubClient := GetGithubClient(content.AccessToken, communicator)
		userLogin, err := GetGithubUserLogin(githubClient)
		if err != nil {
			communicator.SetSend("logout", "Can't retrieve the authenticated user.")
			communicator.Finished(FileContentFinished)
			return
		}
		content.Body, err = GetGithubFileContent(githubClient, userLogin, content.RepositoryName, "content/"+content.Title)
		if err != nil {
			communicator.SetSend("error", "Can't retrieve the file content.")
			communicator.Finished(FileContentFinished)
			return
		}
		communicator.SetSend("content set", content)
		communicator.Finished(FileContentFinished)
	}()
}

// UpdateContent updates the content of a github file.
//
// Happy Path:
// 1. Decode JSON
// 2. Get the github authenticated user
// 3. Update the cotent of a the file
// 4. Send Success message to the clien
func UpdateContent(communicator Communicator, data interface{}) {
	var content Content
	err := mapstructure.Decode(data, &content)
	if err != nil {
		communicator.SetSend("error", "Error decoding json:"+err.Error())
		return
	}
	communicator.NewFinishedChannel(PublishContentFinished)
	go func() {
		githubClient := GetGithubClient(content.AccessToken, communicator)
		userLogin, err := GetGithubUserLogin(githubClient)
		if err != nil {
			communicator.SetSend("logout", "Can't retrieve the authenticated user.")
			communicator.Finished(PublishContentFinished)
			return
		}
		githubFileContentOpt := GetFileContentOptions(content.CommitMessage, content.Branch, userLogin)
		err = UpdateGithubFileContent(githubClient, githubFileContentOpt, content.RepositoryName, "content/"+content.Title)
		if err != nil {
			communicator.SetSend("error", "Unnable to update the content.")
			communicator.Finished(PublishContentFinished)
			return
		}
		// TODO: Send the real shiat to the client.
	}()
}

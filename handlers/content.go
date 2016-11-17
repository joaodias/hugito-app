package handlers

import (
	"github.com/mitchellh/mapstructure"
)

// ContentList represents all the content entries related to a given repository
// and branch.
type ContentList struct {
	Name        string   `json:"name"`
	Branch      string   `json:"branch"`
	Titles      []string `json:"title"`
	Path        string   `json:"path"`
	AccessToken string   `json:"accessToken"`
}

// Content represents the content of a repository and branch. It refers to a content entry
// in the repository.
type Content struct {
	RepositoryName string `json:"repositoryName"`
	Branch         string `json:"branch"`
	Title          string `json:"title"`
	Path           string `json:"path"`
	Body           string `json:"content"`
	Commit         `json:"commit"`
	AccessToken    string `json:"accessToken"`
}

// GetContentList gets the repository content list.
func GetContentList(communicator Communicator, data interface{}) {
	var contentList ContentList
	err := mapstructure.Decode(data, &contentList)
	if err != nil {
		communicator.SetSend("error", "Error decoding json:"+err.Error())
		return
	}
	githubClient := GetGithubClient(contentList.AccessToken, communicator)
	userLogin, err := GetGithubUserLogin(githubClient)
	if err != nil {
		communicator.SetSend("logout", "Can't retrieve the authenticated user.")
		return
	}
	githubContentGetOpt := GetRepositoryContentGetOptions(contentList.Branch)
	contentList.Titles, err = GetGithubRepositoryTree(githubClient, userLogin, githubContentGetOpt, contentList.Name, contentList.Path)
	if err != nil {
		communicator.SetSend("error", "Can't retrieve the content list.")
		return
	}
	communicator.SetSend("content list", contentList)
}

// CreateContent creates a new github content file.
func CreateContent(communicator Communicator, data interface{}) {
	var content Content
	err := mapstructure.Decode(data, &content)
	if err != nil {
		communicator.SetSend("error", "Error decoding json:"+err.Error())
		return
	}
	githubClient := GetGithubClient(content.AccessToken, communicator)
	user, err := GetGithubUser(githubClient)
	if err != nil {
		communicator.SetSend("logout", "Can't retrieve the authenticated user.")
		return
	}
	githubFileContentOpt := GetFileContentOptions(user, content)
	content.Commit, err = CreateGithubFileContent(githubClient, githubFileContentOpt, content.RepositoryName, content.Path+"/"+content.Title)
	if err != nil {
		communicator.SetSend("error", "Unnable to create the content.")
		return
	}
	communicator.SetSend("content create", content)
}

// GetFileContent gets the content of a file from the github repository.
func GetFileContent(communicator Communicator, data interface{}) {
	var content Content
	err := mapstructure.Decode(data, &content)
	if err != nil {
		communicator.SetSend("error", "Error decoding json:"+err.Error())
		return
	}
	githubClient := GetGithubClient(content.AccessToken, communicator)
	userLogin, err := GetGithubUserLogin(githubClient)
	if err != nil {
		communicator.SetSend("logout", "Can't retrieve the authenticated user.")
		return
	}
	githubContentGetOpt := GetRepositoryContentGetOptions(content.Branch)
	content.Body, err = GetGithubFileContent(githubClient, githubContentGetOpt, userLogin, content.RepositoryName, content.Path+"/"+content.Title)
	if err != nil {
		communicator.SetSend("error", "Can't retrieve the file content.")
		return
	}
	communicator.SetSend("content set", content)
}

// UpdateContent updates the content of a github file.
func UpdateContent(communicator Communicator, data interface{}) {
	var content Content
	err := mapstructure.Decode(data, &content)
	if err != nil {
		communicator.SetSend("error", "Error decoding json:"+err.Error())
		return
	}
	githubClient := GetGithubClient(content.AccessToken, communicator)
	user, err := GetGithubUser(githubClient)
	if err != nil {
		communicator.SetSend("logout", "Can't retrieve the authenticated user.")
		return
	}
	githubContentGetOpt := GetRepositoryContentGetOptions(content.Branch)
	githubFileContentOpt := GetFileContentOptions(user, content)
	*githubFileContentOpt.SHA, err = GetGithubFileSHA(githubClient, githubContentGetOpt, user.Login, content.RepositoryName, content.Path+"/"+content.Title)
	if err != nil {
		communicator.SetSend("error", "Unnable to get content information.")
		return
	}
	err = UpdateGithubFileContent(githubClient, githubFileContentOpt, content.RepositoryName, content.Path+"/"+content.Title)
	if err != nil {
		communicator.SetSend("error", "Unnable to update the content.")
		return
	}
	communicator.SetSend("content update", "Content Successfully Published.")
}

// RemoveContent removes an already existent github content file.
func RemoveContent(communicator Communicator, data interface{}) {
	var content Content
	err := mapstructure.Decode(data, &content)
	if err != nil {
		communicator.SetSend("error", "Error decoding json:"+err.Error())
		return
	}
	githubClient := GetGithubClient(content.AccessToken, communicator)
	user, err := GetGithubUser(githubClient)
	if err != nil {
		communicator.SetSend("logout", "Can't retrieve the authenticated user.")
		return
	}
	githubContentGetOpt := GetRepositoryContentGetOptions(content.Branch)
	githubFileContentOpt := GetFileContentOptions(user, content)
	*githubFileContentOpt.SHA, err = GetGithubFileSHA(githubClient, githubContentGetOpt, user.Login, content.RepositoryName, content.Path+"/"+content.Title)
	if err != nil {
		communicator.SetSend("error", "Unnable to get content information.")
		return
	}
	err = RemoveGithubFileContent(githubClient, githubFileContentOpt, content.RepositoryName, content.Path+"/"+content.Title)
	if err != nil {
		communicator.SetSend("error", "Unnable to remove the content.")
		return
	}
	content.Body = ""
	communicator.SetSend("content remove", content)
}

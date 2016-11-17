package handlers

import (
	"github.com/mitchellh/mapstructure"
)

// Repositories represents the repository parameters exchanged between the
// server and the client.
type Repositories struct {
	Names       []string `json:"names"`
	AccessToken string   `json:"accessToken"`
}

// Repository represents a single repository exchanded between the client and
// the server.
type Repository struct {
	Name        string `json:"name"`
	Branch      string `json:"branch"`
	AccessToken string `json:"accessToken"`
}

// GetRepository gets the repositories for an authorized user.
func GetRepository(communicator Communicator, data interface{}) {
	var repositories Repositories
	err := mapstructure.Decode(data, &repositories)
	if err != nil {
		communicator.SetSend("error", "Error decoding json:"+err.Error())
		return
	}
	githubClient := GetGithubClient(repositories.AccessToken, communicator)
	repositories.Names, err = GetGithubRepositories(githubClient)
	if err != nil {
		communicator.SetSend("logout", "Cannot get the user repositories.")
		return
	}
	communicator.SetSend("repositories set", repositories)
}

// ValidateRepository checks if a given repository has a valid hugo
// configuration in the master branch.
func ValidateRepository(communicator Communicator, data interface{}) {
	var repository Repository
	err := mapstructure.Decode(data, &repository)
	if err != nil {
		communicator.SetSend("error", "Error decoding json:"+err.Error())
		return
	}
	githubClient := GetGithubClient(repository.AccessToken, communicator)
	userLogin, err := GetGithubUserLogin(githubClient)
	if err != nil {
		communicator.SetSend("logout", "Can't retrieve the authenticated user.")
		return
	}
	githubContentGetOpt := GetRepositoryContentGetOptions(repository.Branch)
	repositoryTree, err := GetGithubRepositoryTree(githubClient, userLogin, githubContentGetOpt, repository.Name, "")
	if err != nil {
		communicator.SetSend("error", "Can't retrieve selected repository.")
		return
	}
	isValid := IsGithubRepositoryValid(repositoryTree)
	if !isValid {
		communicator.SetSend("error", "Repository is not valid.")
		return
	}
	communicator.SetSend("repository validate", "Repository is valid.")
}

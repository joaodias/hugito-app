package handlers

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
)

// Repositories represents the repository parameters exchanged between the server and the client.
type Repositories struct {
	Names       []string `json:"names"`
	AccessToken string   `json:"accessToken"`
}

// GetRepository gets the repositories for an authorized user.
//
// Happy path:
// 1. Decode JSON
// 3. Get the repositories for the user
// 4. Send repositories to the client
func GetRepository(communicator Communicator, data interface{}) {
	var repositories Repositories
	err := mapstructure.Decode(data, &repositories)
	if err != nil {
		communicator.SetSend("error", "Error decoding json:"+err.Error())
		return
	}
	communicator.NewFinishedChannel(RepositoryFinished)
	go func() {
		githubClient := GetGithubClient(repositories.AccessToken, communicator)
		repositories.Names, err = GetGithubRepositories(githubClient)
		if err != nil {
			communicator.SetSend("logout", "Cannot get the user repositories.")
			communicator.Finished(RepositoryFinished)
			return
		}
		communicator.SetSend("repository set", repositories)
		communicator.Finished(RepositoryFinished)
	}()
}

func UnsubscribeRepository(communicator Communicator, data interface{}) {
	fmt.Print("Repositories unsubscribe \n")
}

func AddRepository(communicator Communicator, data interface{}) {
	fmt.Print("Add repository \n")
}

func RemoveRepository(communicator Communicator, data interface{}) {
	fmt.Print("Remove repository \n")
}

// ValidateRepository checks if a given repository has a valid hugo
// configuration in the master branch.
//
// Happy Path:
// 1. Decode JSON
// 2. Get the github authenticated user
// 3. Get the repository file structure from github
// 4. Check if the given file structure contains the reference file structure.
// 5. Send a repository validate message to the client
func ValidateRepository(communicator Communicator, data interface{}) {
	type Repository struct {
		Name        string `json:"name"`
		AccessToken string `json:"accessToken"`
	}
	var repository Repository
	err := mapstructure.Decode(data, &repository)
	if err != nil {
		communicator.SetSend("error", "Error decoding json:"+err.Error())
		return
	}
	communicator.NewFinishedChannel(RepositoryFinished)
	go func() {
		githubClient := GetGithubClient(repository.AccessToken, communicator)
		userLogin, err := GetGithubUserLogin(githubClient)
		if err != nil {
			communicator.SetSend("logout", "Can't retrieve the authenticated user.")
			communicator.Finished(RepositoryFinished)
			return
		}
		repositoryTree, err := GetGithubRepositoryTree(githubClient, userLogin, repository.Name)
		if err != nil {
			communicator.SetSend("error", "Can't retrieve selected repository.")
			communicator.Finished(RepositoryFinished)
			return
		}
		isValid := IsGithubRepositoryValid(repositoryTree)
		if !isValid {
			communicator.SetSend("error", "Repository is not valid.")
			communicator.Finished(RepositoryFinished)
			return
		}
		communicator.SetSend("repository validate", "Repository is valid.")
		communicator.Finished(RepositoryFinished)
	}()
}

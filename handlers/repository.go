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
	communicator.NewFinishedChannel(RepositoryFinished)
	var repositories Repositories
	go func() {
		err := mapstructure.Decode(data, &repositories)
		if err != nil {
			communicator.SetSend("error", "Error decoding json:"+err.Error())
			communicator.Finished(RepositoryFinished)
			return
		}
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

func ValidateRepository(communicator Communicator, data interface{}) {
	fmt.Print("Validate repository \n")
}

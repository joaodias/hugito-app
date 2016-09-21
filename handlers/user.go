package handlers

import (
	"github.com/mitchellh/mapstructure"
)

// User represents the user information exchanged between the server and the client.
type User struct {
	Name        string `json:"name"`
	AccessToken string `json:"accessToken"`
}

// GetUser gets an user.
//
// The happy flow:
// 1. Decode received data
// 3. Get the github client for this user
// 4. Get the github UserName
func GetUser(communicator Communicator, data interface{}) {
	communicator.NewFinishedChannel(UserFinished)
	var user User
	go func() {
		err := mapstructure.Decode(data, &user)
		if err != nil {
			communicator.SetSend("error", "Error decoding json: ")
			communicator.Finished(UserFinished)
			return
		}
		githubClient := GetGithubClient(user.AccessToken, communicator)
		user.Name, err = GetGithubUserName(githubClient)
		if err != nil {
			communicator.SetSend("logout", "Cannot get the authorized user.")
			communicator.Finished(UserFinished)
			return
		}
		communicator.SetSend("user set", user)
		communicator.Finished(UserFinished)
	}()
}

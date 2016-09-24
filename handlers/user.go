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
	var user User
	err := mapstructure.Decode(data, &user)
	if err != nil {
		communicator.SetSend("error", "Error decoding json: ")
		return
	}
	communicator.NewFinishedChannel(UserFinished)
	go func() {
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

package handlers

import (
	"github.com/mitchellh/mapstructure"
)

// User represents the user information exchanged between the server and the client.
type User struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	Login       string `json:"login"`
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
	githubClient := GetGithubClient(user.AccessToken, communicator)
	user, err = GetGithubUser(githubClient)
	if err != nil {
		communicator.SetSend("logout", "Cannot get the authorized user.")
		return
	}
	communicator.SetSend("user set", user)
}

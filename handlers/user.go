package handlers

import (
	"github.com/mitchellh/mapstructure"
	"time"
)

// User represents the user information exchanged between the server and the client.
type User struct {
	CreatedAt   time.Time `json:"createdAt" gorethink:"createdAt"`
	Name        string    `json:"name" gorethink:"name"`
	Email       string    `json:"email" gorethink:"email"`
	Login       string    `json:"login" gorethink:"login"`
	AccessToken string    `json:"accessToken" gorethink:"accessToken"`
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
	authorizedUser, err := GetGithubUser(githubClient)
	if err != nil {
		communicator.SetSend("logout", "Cannot get the authorized user.")
		return
	}
	authorizedUser.AccessToken = user.AccessToken
	authorizedUser.CreatedAt = time.Now()
	err = communicator.GetDBSession().AddUser(authorizedUser)
	if err != nil {
		communicator.SetSend("error", "Could not register the user. Please try again.")
		return
	}
	communicator.SetSend("user set", authorizedUser)
}

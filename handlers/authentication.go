package handlers

import (
	utils "github.com/joaodias/go-codebase"
	"github.com/mitchellh/mapstructure"
	"golang.org/x/oauth2"
)

// Authentication represents the data exchanged between the communicator and the server
type Authentication struct {
	Authenticated string `json:"authenticated"`
	Code          string `json:"code"`
	State         string `json:"state"`
	ReceivedState string `json:"receivedState"`
}

// Authenticate performs the authentication of an user using Github oauth.
func Authenticate(communicator Communicator, data interface{}) {
	var authentication Authentication
	err := mapstructure.Decode(data, &authentication)
	if err != nil {
		communicator.SetSend("error", "Error decoding json: "+err.Error())
		return
	}
	if !isStateValid(authentication.State, authentication.ReceivedState) {
		communicator.SetSend("error", "received state and state are different.")
		return
	}
	accessToken, err := authentication.getToken(communicator.GetOauthConfiguration())
	if err != nil {
		communicator.SetSend("error", "Error getting the access token.")
		return
	}
	communicator.SetSend("authenticated set", accessToken)
}

func (authentication *Authentication) getToken(oauthConf *oauth2.Config) (string, error) {
	token, err := oauthConf.Exchange(oauth2.NoContext, authentication.Code)
	if err != nil {
		return "", err
	}
	return token.AccessToken, nil
}

func isStateValid(state, receivedState string) bool {
	return utils.AreStringsEqual(state, receivedState)
}

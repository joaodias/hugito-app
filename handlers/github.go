package handlers

import (
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"net/http"
)

// GithubWrapper wraps the github methods from the Github API. All the needed mehtods are wrapped for easier mocking
type GithubWrapper interface {
	GetChecker() Checker
	GetNewClienter() NewClienter
}

// Checker is a function type that checks if a token is authorized in the
// application.
type Checker func(string, string) (*github.Authorization, *github.Response, error)

// NewClienter is a function type that creates a new github client.
type NewClienter func(*http.Client) *github.Client

// IsUserAuthenticated checks if the user is authenticated given an access
// token. It contacts the Github api and checks the validity of the given token
// for this application.
func IsUserAuthenticated(accessToken string, checker Checker) bool {
	_, response, _ := checker(ClientID, accessToken)
	if response.StatusCode == 404 {
		return false
	}
	return true
}

// GetGithubClient gets a client given an access token and an oauth
// configuration. The github client is used to make requests to the github api.
func GetGithubClient(accessToken string, communicator Communicator) *github.Client {
	var token = &oauth2.Token{
		AccessToken: accessToken,
	}
	oauthConf := communicator.GetOauthConfiguration()
	oauthClient := oauthConf.Client(oauth2.NoContext, token)
	newClienter := communicator.GetNewClienter()
	return newClienter(oauthClient)
}

// GetGithubUserName gets the github user name respective to the previouslly
// authorized github client.
func GetGithubUserName(githubClient *github.Client) (string, error) {
	user, _, err := githubClient.Users.Get("")
	if err != nil {
		return "", err
	}
	return *user.Name, nil
}

// GetChecker returns an authservice from the github api. Authservice is used
// to check wether an user is authorized or not.
func (socketClient *SocketClient) GetChecker() Checker {
	var authService github.AuthorizationsService
	return authService.Check
}

// GetNewClienter returns a newClient function from the github api. It is a
// wrapper to the api to improve testability and flexibility.
func (socketClient *SocketClient) GetNewClienter() NewClienter {
	return github.NewClient
}

package handlers

import (
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"net/http"
)

// GithubWrapper wraps the github methods from the Github API. All the needed mehtods are wrapped for easier mocking
type GithubWrapper interface {
	GetNewClienter() NewClienter
}

// NewClienter is a function type that creates a new github client.
type NewClienter func(*http.Client) *github.Client

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

// GetGithubRepositories gets the github repositories to the previously
// authorized client.
func GetGithubRepositories(githubClient *github.Client) ([]string, error) {
	opt := &github.RepositoryListOptions{
		Type: "all",
	}
	repos, _, err := githubClient.Repositories.List("", opt)
	if err != nil {
		return nil, err
	}
	repositoriesName := make([]string, len(repos))
	for i := 0; i < len(repos); i++ {
		repositoriesName[i] = *repos[i].Name
	}

	return repositoriesName, nil
}

// GetNewClienter returns a NewClient function from the github api. It is a
// wrapper to the api to improve testability and flexibility.
func (socketClient *SocketClient) GetNewClienter() NewClienter {
	return github.NewClient
}

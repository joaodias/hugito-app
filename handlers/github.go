package handlers

import (
	"github.com/google/go-github/github"
	utils "github.com/joaodias/hugito-app/utils"
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

// GetGithubUserLogin gets the github user login name respective to the previouslly
// authorized github client.
func GetGithubUserLogin(githubClient *github.Client) (string, error) {
	user, _, err := githubClient.Users.Get("")
	if err != nil {
		return "", err
	}
	return *user.Login, nil
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

// GetGithubRepositoryTree gets the tree of files in the given path of a
// repository of a given user.
func GetGithubRepositoryTree(githubClient *github.Client, userLogin string, repositoryName string, path string) ([]string, error) {
	opt := &github.RepositoryContentGetOptions{}
	_, githubRepositoryTree, _, err := githubClient.Repositories.GetContents(userLogin, repositoryName, path, opt)
	if err != nil {
		return []string{}, err
	}
	repositoryTree := make([]string, len(githubRepositoryTree))
	for i := 0; i < len(githubRepositoryTree); i++ {
		repositoryTree[i] = *githubRepositoryTree[i].Name
	}
	return repositoryTree, nil
}

// IsGithubRepositoryValid checks if a given repository tree matches the
// criteria to be a valid HUGO repository.
func IsGithubRepositoryValid(repositoryTree []string) bool {
	referenceTree := []string{"content", "config.toml", "public", "themes"}
	isValid := utils.ContainsSubArray(repositoryTree, referenceTree)
	return isValid
}

// GetNewClienter returns a NewClient function from the github api. It is a
// wrapper to the api to improve testability and flexibility.
func (socketClient *SocketClient) GetNewClienter() NewClienter {
	return github.NewClient
}

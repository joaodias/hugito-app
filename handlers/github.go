package handlers

import (
	"encoding/base64"
	"github.com/google/go-github/github"
	utils "github.com/joaodias/hugito-app/utils"
	"golang.org/x/oauth2"
	"net/http"
)

// Default values to use with the github wrapper.
const (
	DefaultCommitMessage = "Updated by Hugito"
	DefaultBranch        = "master"
	DefaultAuthor        = "Hugito"
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

// GetGithubFileContent gets the content of a HUGO content file.
func GetGithubFileContent(githubClient *github.Client, userLogin string, repositoryName string, path string) (string, error) {
	opt := &github.RepositoryContentGetOptions{}
	fileContent, _, _, err := githubClient.Repositories.GetContents(userLogin, repositoryName, path, opt)
	if err != nil {
		return "", err
	}
	decodedContent, err := base64.StdEncoding.DecodeString(*fileContent.Content)
	if err != nil {
		return "", err
	}
	return string(decodedContent), nil
}

// UpdateGithubFileContent updates the content of an already existent file in
// Github. A github content file object stucture is passed to improve
// testability, flexibility and also to improve the readability of the method.
func UpdateGithubFileContent(githubClient *github.Client, opt *github.RepositoryContentFileOptions, repositoryName string, path string) error {
	_, _, err := githubClient.Repositories.UpdateFile(*opt.Author.Login, repositoryName, path, opt)
	if err != nil {
		return err
	}
	return nil
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

// GetFileContentOptions builds a file content option structure given a set of
// parameters. SHA is ignored and Commiter is interpreted as the same role as
// Author. For the time being The CommitAuthor just consists of the Login field.
func GetFileContentOptions(message string, branch string, author string) *github.RepositoryContentFileOptions {
	if message == "" {
		message = DefaultCommitMessage
	}
	if author == "" {
		author = DefaultAuthor
	}
	if branch == "" {
		branch = DefaultBranch
	}
	return &github.RepositoryContentFileOptions{
		Message: &message,
		Branch:  &branch,
		Author: &github.CommitAuthor{
			Login: &author,
		},
		Committer: &github.CommitAuthor{
			Login: &author,
		},
	}
}

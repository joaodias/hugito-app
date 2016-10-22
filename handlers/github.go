package handlers

import (
	"encoding/base64"
	"github.com/google/go-github/github"
	utils "github.com/joaodias/go-codebase"
	"golang.org/x/oauth2"
	"net/http"
)

// Default values for optional or missing fields.
const (
	DefaultCommitMessage = "Job done by the awesome Hugito"
	DefaultBranch        = "master"
	DefaultContentBody   = "`Hugito created this piece of content just for you <3`"
)

// Commit wrapps the Commit structure from the github api. This structure is represented just by the essential commit information.
type Commit struct {
	SHA     string `json:"sha"`
	Message string `json:"commitMessage"`
	URL     string `json:"url"`
	Author  `json:"author"`
}

// Author is the author of a content entry.
type Author struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Login string `json:"login"`
}

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

// GetGithubUser gets the github user.
func GetGithubUser(githubClient *github.Client) (User, error) {
	user, _, err := githubClient.Users.Get("")
	if err != nil {
		return User{}, err
	}
	return User{
		Name:  *user.Name,
		Email: *user.Email,
		Login: *user.Login,
	}, nil
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
func GetGithubRepositoryTree(githubClient *github.Client, userLogin string, opt *github.RepositoryContentGetOptions, repositoryName string, path string) ([]string, error) {
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

// CreateGithubFileContent creates a new Github file content at a given
// repository and in a given branch.
func CreateGithubFileContent(githubClient *github.Client, opt *github.RepositoryContentFileOptions, repositoryName string, path string) (Commit, error) {
	repositoryContentResponse, _, err := githubClient.Repositories.CreateFile(*opt.Author.Login, repositoryName, path, opt)
	if err != nil {
		return Commit{}, err
	}
	return Commit{
		SHA: *repositoryContentResponse.SHA,
		Author: Author{
			Name:  *repositoryContentResponse.Author.Name,
			Email: *repositoryContentResponse.Author.Email,
		},
	}, nil
}

// GetGithubFileContent gets the content of a HUGO content file.
func GetGithubFileContent(githubClient *github.Client, opt *github.RepositoryContentGetOptions, userLogin string, repositoryName string, path string) (string, error) {
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

// GetGithubFileSHA gets the SHA for the given file.
func GetGithubFileSHA(githubClient *github.Client, opt *github.RepositoryContentGetOptions, userLogin string, repositoryName string, path string) (string, error) {
	content, _, _, err := githubClient.Repositories.GetContents(userLogin, repositoryName, path, opt)
	if err != nil {
		return "", err
	}
	return *content.SHA, nil
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

// RemoveGithubFileContent removes an already existing Github file content at a
// given repository and in a given branch.
func RemoveGithubFileContent(githubClient *github.Client, opt *github.RepositoryContentFileOptions, repositoryName string, path string) error {
	_, _, err := githubClient.Repositories.DeleteFile(*opt.Author.Login, repositoryName, path, opt)
	if err != nil {
		return err
	}
	return nil
}

// IsGithubRepositoryValid checks if a given repository tree matches the
// criteria to be a valid HUGO repository.
func IsGithubRepositoryValid(repositoryTree []string) bool {
	referenceTree := []string{"config.toml", "public", "themes"}
	isValid := utils.ContainsSubArray(repositoryTree, referenceTree)
	return isValid
}

// GetNewClienter returns a NewClient function from the github api. It is a
// wrapper to the api to improve testability and flexibility.
func (socketClient *SocketClient) GetNewClienter() NewClienter {
	return github.NewClient
}

// GetRepositoryContentGetOptions gets the github repositoy content get options based on a given branch.
func GetRepositoryContentGetOptions(branch string) *github.RepositoryContentGetOptions {
	if branch == "" {
		branch = DefaultBranch
	}
	return &github.RepositoryContentGetOptions{
		Ref: branch,
	}
}

// GetFileContentOptions builds a file content option structure given an user
// and the content. Commiter is interpreted as the same role as Author.
func GetFileContentOptions(user User, content Content) *github.RepositoryContentFileOptions {
	if content.Commit.Message == "" {
		content.Commit.Message = DefaultCommitMessage
	}
	if content.Branch == "" {
		content.Branch = DefaultBranch
	}
	if content.Body == "" {
		content.Body = DefaultContentBody
	}
	return &github.RepositoryContentFileOptions{
		Message: &content.Commit.Message,
		Branch:  &content.Branch,
		Content: []byte(content.Body),
		SHA:     &content.SHA,
		Author: &github.CommitAuthor{
			Login: &user.Login,
			Email: &user.Email,
			Name:  &user.Name,
		},
		Committer: &github.CommitAuthor{
			Login: &user.Login,
			Email: &user.Email,
			Name:  &user.Name,
		},
	}
}

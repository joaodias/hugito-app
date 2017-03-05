package infrastructure

import (
	"encoding/base64"
	"errors"
	"github.com/google/go-github/github"
	utils "github.com/joaodias/go-codebase/files"
	"github.com/joaodias/go-codebase/strings"
	"github.com/joaodias/hugito-backend/domain"
	"golang.org/x/oauth2"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"time"
)

const (
	projectBranchRootPath = ""
	defaultBody           = "`Hugito created this piece of content just for you <3`"
	compression           = "zipball"
	compressionExtension  = ".zip"
	extractionFolder      = "./"
)

var (
	defaultCommitMessage = "Commited with <3 by the awesome Hugito."
)

type Github struct {
	Client             *github.Client
	User               *github.User
	Data               interface{}
	OauthConfiguration *oauth2.Config
}

// Tree is a custom tree that extends the github tree.
type Tree struct {
	Name string
	// Parent refers to the parent tree of a tree. If a tree has no parent then
	// it is equal to the unique identied parent folder containing the unziped
	// repository.
	Parent string
	// Level represents that deepness of a tree. The deeper it is, the deper the folder level.
	Level      int
	GithubTree *github.Tree
}

// GetUser gets the github user
func (gh *Github) GetUser(accessToken string, oauthConfiguration *oauth2.Config) (*domain.User, error) {
	client := gh.newGithubClient(accessToken, oauthConfiguration)
	gh.Client = client
	user, err := gh.getGithubUser()
	if err != nil {
		return nil, err
	}
	return &domain.User{
		Name:  *user.Name,
		Email: *user.Email,
		Login: *user.Login,
	}, nil
}

// Validates if a given repository is a valid repository. To be valid a repository should have some specific files or
// folders.
func (gh *Github) ValidateRepository(accessToken string, oauthConfiguration *oauth2.Config, repository domain.Repository) (bool, error) {
	client := gh.newGithubClient(accessToken, oauthConfiguration)
	gh.Client = client
	user, err := gh.getGithubUser()
	if err != nil {
		return false, err
	}
	gh.User = user
	repositoryTree, err := gh.getGithubRepositoryTree(repository.Name, repository.ProjectBranch, projectBranchRootPath)
	if err != nil {
		return false, err
	}
	referenceTree := []string{"config.toml", "public", "themes"}
	isValid := strings.ContainsSubArray(repositoryTree, referenceTree)
	return isValid, nil
}

// ListContentTitles lists the content of a set of github files.
func (gh *Github) ListContentTitles(content domain.Content, oauthConfiguration *oauth2.Config) ([]string, error) {
	client := gh.newGithubClient(content.AccessToken, oauthConfiguration)
	gh.Client = client
	user, err := gh.getGithubUser()
	if err != nil {
		return nil, err
	}
	gh.User = user
	return gh.getGithubRepositoryTree(content.RepositoryName, content.ProjectBranch, content.Path)
}

// GetFileContent retrieves the content of a github file.
func (gh *Github) GetFileContent(content domain.Content, oauthConfiguration *oauth2.Config) (*string, error) {
	client := gh.newGithubClient(content.AccessToken, oauthConfiguration)
	gh.Client = client
	user, err := gh.getGithubUser()
	if err != nil {
		return nil, err
	}
	gh.User = user
	gh.Data = content
	fileContent, err := gh.getGithubFileContent()
	if err != nil {
		return nil, err
	}
	decodedContent, err := base64.StdEncoding.DecodeString(*fileContent)
	if err != nil {
		return nil, err
	}
	decodedContentString := string(decodedContent)
	return &decodedContentString, nil
}

// UpdateFileContent updates the content of a github file.
func (gh *Github) UpdateFileContent(content domain.Content, oauthConfiguration *oauth2.Config) error {
	client := gh.newGithubClient(content.AccessToken, oauthConfiguration)
	gh.Client = client
	user, err := gh.getGithubUser()
	if err != nil {
		return err
	}
	gh.User = user
	gh.Data = content
	contentFileOpt, err := gh.buildContentFileOpt(content.ProjectBranch)
	if err != nil {
		return err
	}
	_, _, err = client.Repositories.UpdateFile(*contentFileOpt.Author.Login, content.RepositoryName, content.Path+"/"+content.Title, contentFileOpt)
	if err != nil {
		return err
	}
	return nil
}

// CreateContentFile creates a content file in github with the given content information.
func (gh *Github) CreateContentFile(content domain.Content, oauthConfiguration *oauth2.Config) (*domain.Content, error) {
	client := gh.newGithubClient(content.AccessToken, oauthConfiguration)
	gh.Client = client
	user, err := gh.getGithubUser()
	if err != nil {
		return nil, err
	}
	gh.User = user
	gh.Data = content
	commit, err := gh.createGithubFileContent()
	if err != nil {
		return nil, err
	}
	content.Commit = *commit
	return &content, nil
}

// RemoveContentFile removes an already existing Github file content at a
// given repository and in a given branch.
func (gh *Github) RemoveContentFile(content domain.Content, oauthConfiguration *oauth2.Config) error {
	client := gh.newGithubClient(content.AccessToken, oauthConfiguration)
	gh.Client = client
	user, err := gh.getGithubUser()
	if err != nil {
		return err
	}
	gh.User = user
	gh.Data = content
	contentFileOpt, err := gh.buildContentFileOpt(content.ProjectBranch)
	if err != nil {
		return err
	}
	_, _, err = client.Repositories.DeleteFile(*contentFileOpt.Author.Login, content.RepositoryName, content.Path+"/"+content.Title, contentFileOpt)
	if err != nil {
		return err
	}
	return nil
}

// DownloadGithubContents downloads the github contents from the HUGO project branch.
func (gh *Github) DownloadContents(content domain.Content, oauthConfiguration *oauth2.Config) (*string, error) {
	client := gh.newGithubClient(content.AccessToken, oauthConfiguration)
	user, err := gh.getGithubUser()
	if err != nil {
		return nil, err
	}
	contentGetOpt := &github.RepositoryContentGetOptions{
		Ref: content.ProjectBranch,
	}
	url, response, err := client.Repositories.GetArchiveLink(*user.Login, content.RepositoryName, compression, contentGetOpt)
	if err != nil || response.StatusCode != 302 {
		return nil, errors.New("Can't get Github repository.")
	}
	sourcePath, err := utils.DownloadToFile(utils.Client{HTTP: &utils.RealWebClient{}}, url.String(), extractionFolder, content.RepositoryName+compressionExtension, true)
	if err != nil {
		return nil, err
	}
	return &sourcePath, nil
}

// RemoveDonwloadedContents removes the path of downloaded contents from github.
func (gh *Github) RemoveDownloadedContents(path string) error {
	return os.RemoveAll(path)
}

// PushFiles pushes files from a source path to github.
func (gh *Github) PushFiles(content domain.Content, oauthConfiguration *oauth2.Config, sourcePath string) error {
	client := gh.newGithubClient(content.AccessToken, oauthConfiguration)
	gh.Client = client
	user, err := gh.getGithubUser()
	if err != nil {
		return err
	}
	gh.User = user
	gh.Data = content
	treeMap, err := gh.buildGithubTrees(sourcePath)
	if err != nil {
		return err
	}
	mainTree, err := gh.publishTrees(treeMap, utils.GetPathRoot(sourcePath))
	if err != nil {
		return err
	}
	commitSHA, commitURL, err := gh.createCommit(mainTree)
	if err != nil {
		return err
	}
	err = gh.updateBranchHead(&content.PublicBranch, commitSHA, commitURL)
	if err != nil {
		return err
	}
	return err
}

func (gh *Github) newGithubClient(accessToken string, oauthConfiguration *oauth2.Config) *github.Client {
	var token = &oauth2.Token{
		AccessToken: accessToken,
	}
	oauthClient := oauthConfiguration.Client(oauth2.NoContext, token)
	return github.NewClient(oauthClient)
}

func (gh *Github) getGithubUser() (*github.User, error) {
	client := gh.Client
	user, _, err := client.Users.Get("")
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (gh *Github) getGithubRepositoryTree(repositoryName, projectBranch, root string) ([]string, error) {
	client := gh.Client
	user := gh.User
	opt := &github.RepositoryContentGetOptions{
		Ref: projectBranch,
	}
	_, githubRepositoryTree, _, err := client.Repositories.GetContents(*user.Login, repositoryName, root, opt)
	if err != nil {
		return []string{}, err
	}
	repositoryTree := make([]string, len(githubRepositoryTree))
	for i := 0; i < len(githubRepositoryTree); i++ {
		repositoryTree[i] = *githubRepositoryTree[i].Name
	}
	return repositoryTree, nil
}

func (gh *Github) getGithubFileContent() (*string, error) {
	client := gh.Client
	user := gh.User
	content := gh.Data.(domain.Content)
	opt := &github.RepositoryContentGetOptions{
		Ref: content.ProjectBranch,
	}
	fileContent, _, _, err := client.Repositories.GetContents(*user.Login, content.RepositoryName, content.Path+"/"+content.Title, opt)
	if err != nil {
		return nil, err
	}
	return fileContent.Content, nil
}

func (gh *Github) buildContentFileOpt(branch string) (*github.RepositoryContentFileOptions, error) {
	client := gh.Client
	user := gh.User
	content := gh.Data.(domain.Content)
	contentGetOpt := &github.RepositoryContentGetOptions{
		Ref: branch,
	}
	repositoryContent, _, _, err := client.Repositories.GetContents(*user.Login, content.RepositoryName, content.Path+"/"+content.Title, contentGetOpt)
	if err != nil {
		return nil, err
	}
	if content.Body == "" {
		content.Body = defaultBody
	}
	contentFileOpt := &github.RepositoryContentFileOptions{
		Message: &defaultCommitMessage,
		Branch:  &branch,
		Content: []byte(content.Body),
		SHA:     repositoryContent.SHA,
		Author: &github.CommitAuthor{
			Login: user.Login,
			Email: user.Email,
			Name:  user.Name,
		},
		Committer: &github.CommitAuthor{
			Login: user.Login,
			Email: user.Email,
			Name:  user.Name,
		},
	}
	return contentFileOpt, nil
}

func (gh *Github) createGithubFileContent() (*domain.Commit, error) {
	client := gh.Client
	user := gh.User
	content := gh.Data.(domain.Content)
	if content.Body == "" {
		content.Body = defaultBody
	}
	opt := &github.RepositoryContentFileOptions{
		Message: &defaultCommitMessage,
		Branch:  &content.ProjectBranch,
		Content: []byte(content.Body),
		SHA:     &content.Commit.SHA,
		Author: &github.CommitAuthor{
			Login: user.Login,
			Email: user.Email,
			Name:  user.Name,
		},
		Committer: &github.CommitAuthor{
			Login: user.Login,
			Email: user.Email,
			Name:  user.Name,
		},
	}
	repositoryContentResponse, _, err := client.Repositories.CreateFile(*opt.Author.Login, content.RepositoryName, content.Path+"/"+content.Title, opt)
	if err != nil {
		return nil, err
	}
	return &domain.Commit{
		SHA:   *repositoryContentResponse.SHA,
		Name:  *repositoryContentResponse.Author.Name,
		Email: *repositoryContentResponse.Author.Email,
	}, nil
}

func (gh *Github) buildGithubTrees(sourcePath string) (map[string]*Tree, error) {
	treeMap := make(map[string]*Tree)
	baseDirectory := sourcePath
	parseFileInfo := func(filePath string, fileInfo os.FileInfo, err error) error {
		if filePath == baseDirectory {
			return nil
		}
		fileDirectory := path.Dir(filePath)
		if treeMap[fileDirectory] == nil {
			treeMap[fileDirectory], err = gh.newTreeMapEntry(filePath)
		} else {
			treeEntry, err := gh.newGithubTreeEntry(filePath)
			if err != nil {
				return err
			}
			treeMap[fileDirectory].GithubTree.Entries = append(treeMap[fileDirectory].GithubTree.Entries, *treeEntry)
		}
		return nil
	}
	err := filepath.Walk(sourcePath, parseFileInfo)
	if err != nil {
		return nil, err
	}
	return treeMap, nil
}

func (gh *Github) createCommit(tree *github.Tree) (*string, *string, error) {
	client := gh.Client
	user := gh.User
	content := gh.Data.(domain.Content)
	now := time.Now()
	author := &github.CommitAuthor{
		Date:  &now,
		Name:  user.Name,
		Email: user.Email,
		Login: user.Login,
	}
	responseCommit, _, err := client.Git.CreateCommit(*user.Login, content.RepositoryName, &github.Commit{
		Author:    author,
		Committer: author,
		Message:   &defaultCommitMessage,
		Tree:      tree,
	})
	if err != nil {
		return nil, nil, err
	}
	return responseCommit.SHA, responseCommit.URL, err
}

func (gh *Github) updateBranchHead(branch *string, commitSHA *string, commitURL *string) error {
	client := gh.Client
	user := gh.User
	content := gh.Data.(domain.Content)
	refHead := "heads/" + content.PublicBranch
	_, _, err := client.Git.UpdateRef(*user.Login, content.RepositoryName, &github.Reference{
		Ref: &refHead,
		Object: &github.GitObject{
			SHA: commitSHA,
			URL: commitURL,
		}}, true)
	if err != nil {
		return err
	}
	return nil
}

func (gh *Github) getTreeLevel(path string) int {
	currentLevel := 0
	currentPath := path
	for currentPath != "." {
		currentPath = filepath.Dir(currentPath)
		currentLevel++
	}
	return currentLevel
}

func (gh *Github) getMaxTreeLevel(treeMap map[string]*Tree) int {
	currentLevel := 0
	for _, tree := range treeMap {
		if tree.Level > currentLevel {
			currentLevel = tree.Level
		}
	}
	return currentLevel
}

func (gh *Github) publishTrees(treeMap map[string]*Tree, rootPath string) (*github.Tree, error) {
	currentLevel := gh.getMaxTreeLevel(treeMap)
	var mainTree *github.Tree
	for currentLevel != 0 {
		for _, tree := range treeMap {
			if tree.Level == currentLevel {
				responseTree, err := gh.publishTree(tree)
				if err != nil {
					return nil, err
				}
				// In order to publish the parent tree, it should know the sha
				// of all children trees. All of the shas should be added to the
				// parent tree. The main tree dow not need to add its sha to a
				// parent tree. It just needs to return it.
				if tree.Parent != rootPath {
					treeMap = gh.addTreeSHAToParent(responseTree.SHA, tree, treeMap)
				} else {
					mainTree = responseTree
				}
			}
		}
		currentLevel--
	}
	return mainTree, nil
}

func (gh *Github) addTreeSHAToParent(treeSHA *string, tree *Tree, treeMap map[string]*Tree) map[string]*Tree {
	entries := treeMap[tree.Parent].GithubTree.Entries
	for i, githubTree := range entries {
		if *githubTree.Path == tree.Name {
			treeMap[tree.Parent].GithubTree.Entries[i].SHA = treeSHA
		}
	}
	return treeMap
}

func (gh *Github) publishTree(tree *Tree) (*github.Tree, error) {
	content := gh.Data.(domain.Content)
	client := gh.Client
	user := gh.User
	responseTree, _, err := client.Git.CreateTree(*user.Login, content.RepositoryName, "", tree.GithubTree.Entries)
	if err != nil {
		return nil, err
	}
	return responseTree, nil
}

func (gh *Github) newTreeMapEntry(filePath string) (*Tree, error) {
	githubTree := &github.Tree{}
	treeEntry, err := gh.newGithubTreeEntry(filePath)
	if err != nil {
		return nil, err
	}
	githubTree.Entries = append(githubTree.Entries, *treeEntry)
	// We need to get the directory of the file path because in what comes to
	// the tree structure, its concern is the folders and not the files
	fileDirectory := filepath.Dir(filePath)
	return &Tree{
		Name:       filepath.Base(fileDirectory),
		Parent:     filepath.Dir(fileDirectory),
		Level:      gh.getTreeLevel(fileDirectory),
		GithubTree: githubTree,
	}, nil
}

func (gh *Github) newGithubTreeEntry(path string) (*github.TreeEntry, error) {
	var encodedFileContent string
	var treeEntry *github.TreeEntry
	isDirectory, err := utils.IsDirectory(path)
	if err != nil {
		return nil, err
	}
	if isDirectory {
		blobType := "tree"
		blobMode := "040000"
		basePath := filepath.Base(path)
		treeEntry = &github.TreeEntry{
			Path: &basePath,
			Mode: &blobMode,
			Type: &blobType,
		}
	} else {
		blobType := "blob"
		blobMode := "100644"
		fileContent, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, err
		}
		encodedFileContent = base64.StdEncoding.EncodeToString(fileContent)
		encoding := "base64"
		blobSHA, err := gh.publishBlob(&github.Blob{
			Content:  &encodedFileContent,
			Encoding: &encoding,
		})
		if err != nil {
			return nil, err
		}
		basePath := filepath.Base(path)
		treeEntry = &github.TreeEntry{
			Path: &basePath,
			Mode: &blobMode,
			Type: &blobType,
			SHA:  blobSHA,
		}
	}
	return treeEntry, nil
}

func (gh *Github) publishBlob(blob *github.Blob) (*string, error) {
	client := gh.Client
	user := gh.User
	content := gh.Data.(domain.Content)
	responseBlob, _, err := client.Git.CreateBlob(*user.Login, content.RepositoryName, blob)
	if err != nil {
		return nil, err
	}
	return responseBlob.SHA, nil
}

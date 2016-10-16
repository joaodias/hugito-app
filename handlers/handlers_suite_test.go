package handlers_test

import (
	"github.com/google/go-github/github"
	handlers "github.com/joaodias/hugito-app/handlers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"reflect"
	"testing"
)

func TestHandlers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Handlers Suite")
}

func TestGetFileContentOptions(t *testing.T) {
	defaultBranch := handlers.DefaultBranch
	defaultCommitMessage := handlers.DefaultCommitMessage
	defaultContentBody := handlers.DefaultContentBody
	user := handlers.User{
		Name:  "João Dias",
		Email: "diasjoaoac@gmail.com",
		Login: "joaodias",
	}
	content := handlers.Content{
		Commit: handlers.Commit{
			Message: "Cool Message",
			SHA:     "1234",
		},
		Branch: "cool-branch",
		Body:   "cool content",
	}
	incompleteContent := handlers.Content{
		Commit: handlers.Commit{
			Message: "",
			SHA:     "",
		},
		Branch: "",
		Body:   "",
	}
	type args struct {
		user    handlers.User
		content handlers.Content
	}
	tests := []struct {
		name string
		args args
		want *github.RepositoryContentFileOptions
	}{
		{"With all the fields provided", args{user, content}, &github.RepositoryContentFileOptions{
			Message: &content.Message,
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
		}},
		{"With default fields", args{user, incompleteContent}, &github.RepositoryContentFileOptions{
			Message: &defaultCommitMessage,
			Branch:  &defaultBranch,
			Content: []byte(defaultContentBody),
			SHA:     &incompleteContent.SHA,
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
		}},
	}
	for _, tt := range tests {
		if got := handlers.GetFileContentOptions(tt.args.user, tt.args.content); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. GetFileContentOptions() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestGetRepositoryContentGetOptions(t *testing.T) {
	defaultBranch := handlers.DefaultBranch
	type args struct {
		branch string
	}
	tests := []struct {
		name string
		args args
		want *github.RepositoryContentGetOptions
	}{
		{"With branch name", args{"1-cool-branch"}, &github.RepositoryContentGetOptions{
			Ref: "1-cool-branch",
		}},
		{"With default branch name", args{""}, &github.RepositoryContentGetOptions{
			Ref: defaultBranch,
		}},
	}
	for _, tt := range tests {
		if got := handlers.GetRepositoryContentGetOptions(tt.args.branch); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. GetRepositoryContentGetOptions() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

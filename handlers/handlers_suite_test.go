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

// This method diserves a test for itself. I think it as some sensible logic
// that looks like to be better tested with a table test.
func TestGetFileContentOptions(t *testing.T) {
	DefaultCommitMessage := "Updated by Hugito"
	DefaultBranch := "master"
	DefaultAuthor := "Hugito"
	mockMessage := "Cool message"
	mockBranch := "1-cool-branch"
	mockAuthor := "joaodias"
	type args struct {
		message string
		branch  string
		author  string
	}
	tests := []struct {
		name string
		args args
		want *github.RepositoryContentFileOptions
	}{
		{"With all parameters provided", args{mockMessage, mockBranch, mockAuthor}, &github.RepositoryContentFileOptions{
			Message: &mockMessage,
			Branch:  &mockBranch,
			Author: &github.CommitAuthor{
				Login: &mockAuthor,
			},
			Committer: &github.CommitAuthor{
				Login: &mockAuthor,
			},
		}},
		{"With default paramters", args{"", "", ""}, &github.RepositoryContentFileOptions{
			Message: &DefaultCommitMessage,
			Branch:  &DefaultBranch,
			Author: &github.CommitAuthor{
				Login: &DefaultAuthor,
			},
			Committer: &github.CommitAuthor{
				Login: &DefaultAuthor,
			},
		}},
	}
	for _, tt := range tests {
		if got := handlers.GetFileContentOptions(tt.args.message, tt.args.branch, tt.args.author); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. GetFileContentOptions() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestGetRepositoryContentGetOptions(t *testing.T) {
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
			Ref: "master",
		}},
	}
	for _, tt := range tests {
		if got := handlers.GetRepositoryContentGetOptions(tt.args.branch); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. GetRepositoryContentGetOptions() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

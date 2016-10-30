package mocks

import (
	"github.com/google/go-github/github"
	handlers "github.com/joaodias/hugito-app/handlers"
	models "github.com/joaodias/hugito-app/models"
	"golang.org/x/oauth2"
	"net/http"
	"net/url"
)

// Mocking for the Client.

const mockToken = "90d64460d14870c08c81352a05dedd3465940a7c"

type Message struct {
	Name string      `json:"name"`
	Data interface{} `json:"data"`
}

type Authentication struct {
	Authenticated string `json:"authenticated"`
	Code          string `json:"code"`
	State         string `json:"state"`
	ReceivedState string `json:"receivedState"`
}

type Client struct {
	Name      string
	Data      interface{}
	OauthConf *oauth2.Config
	// Error is used to signal an error in case a method needs it for testing.
	IsError bool
}

var _ handlers.Communicator = (*Client)(nil)

func (t *Client) Read() {
	//Do Nothing
}

func (t *Client) Write() {
	// Do Nothing
}

func (t *Client) SetSend(name string, data interface{}) {
	t.Name = name
	t.Data = data
}

func (t *Client) GetOauthConfiguration() *oauth2.Config {
	return t.OauthConf
}

func (t *Client) GetNewClienter() handlers.NewClienter {
	return t.NewClient
}

func (t *Client) GetDBSession() models.DataStorage {
	// Data storage needs an error cause we have to test the scenarios wether
	// the data storage logic returns or not an erro.
	return &DataStorage{
		IsError: t.IsError,
	}
}

func (t *Client) NewClient(httpClient *http.Client) *github.Client {
	client := github.NewClient(nil)
	oauthConf := t.GetOauthConfiguration()
	testServerURL := oauthConf.Endpoint.TokenURL
	client.BaseURL, _ = url.Parse(testServerURL)
	client.UploadURL, _ = url.Parse(testServerURL)
	return client
}

package interfaces

import (
	"encoding/json"
	"errors"
	"github.com/joaodias/go-codebase/strings"
	"github.com/joaodias/hugito-backend/usecases"
	"golang.org/x/oauth2"
	"io"
	"net/http"
)

// UserInteractor manipulates the user entity
type UserInteractor interface {
	New(string, string, string, string) (*usecases.User, error)
	Read(string, *oauth2.Config) (*usecases.User, error)
}

// ContentInteractor manipulates the content entity
type ContentInteractor interface {
	New(content usecases.Content, oauthConfiguration *oauth2.Config) (*usecases.Content, error)
	Remove(content usecases.Content, oauthConfiguration *oauth2.Config) error
	Update(content usecases.Content, oauthConfiguration *oauth2.Config) (*usecases.Content, error)
	List(content usecases.Content, oauthConfiguration *oauth2.Config) ([]usecases.Content, error)
	Find(content usecases.Content, oauthConfiguration *oauth2.Config) (*usecases.Content, error)
	Publish(content usecases.Content, oauthConfiguration *oauth2.Config) error
}

// RepositoryInteractor manipulates the Repository entity
// TODO: Should be just 1 usecase. In this case Validate.
type RepositoryInteractor interface {
	New(string, string, string, string) (*usecases.Repository, error)
	Validate(string, string, string, *oauth2.Config) (bool, error)
}

// WebHandler handles the request provinient from the web.
type WebHandler struct {
	OauthConfiguration   *oauth2.Config
	UserInteractor       UserInteractor
	ContentInteractor    ContentInteractor
	RepositoryInteractor RepositoryInteractor
	Logger               Logger
}

// AuthWebOptions are the options fot the auth endpoint.
func (wh *WebHandler) AuthWebOptions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "access-control-allow-origin, content-type,  access-control-allow-headers, access-control-allow-methods")
	w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, POST")
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

// ContentWebOptions are the options for the content endpoint.
func (wh *WebHandler) ContentWebOptions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "access-control-allow-origin, content-type,  access-control-allow-headers, access-control-allow-methods")
	w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, POST, GET, PUT, DELETE")
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

// Authenticate authenticates an user. Some part of the authentication process takes place on the client. However, for
// obvious reasons, there is the need to verify the authentication in the backend.
func (wh *WebHandler) Authenticate(w http.ResponseWriter, r *http.Request) {
	wh.Logger.Log("[web_Authenticate]")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	type authentication struct {
		State         string `json:"state"`
		Code          string `json:"code"`
		ReceivedState string `json:"receivedState"`
		AccessToken   string `json:"accessToken"`
	}
	decodedAuth, err := decodeJSON(r.Body, &authentication{})
	if err != nil {
		wh.Logger.Log("[web_Authenticate]" + err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	auth := decodedAuth.(*authentication)
	if !wh.isStateValid(auth.State, auth.ReceivedState) {
		wh.Logger.Log("[web_Authenticate] Authentication with different states.")
		http.Error(w, "States are not equal.", http.StatusNotAcceptable)
		return
	}
	accessToken, err := wh.exchangeToken(auth.Code, wh.OauthConfiguration)
	if err != nil {
		wh.Logger.Log("[web_Authenticate]" + err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response, _ := json.Marshal(&authentication{AccessToken: *accessToken})
	wh.Logger.Log("[web_Authenticate] Authentication successful with token: " + string(response))
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

// ReadUser retrieves the user from an external repository
func (wh *WebHandler) ReadUser(w http.ResponseWriter, r *http.Request) {
	wh.Logger.Log("[web_ReadUser]")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	accessToken := r.FormValue("accessToken")
	if accessToken == "" {
		wh.Logger.Log("[web_ReadUser] Missing access token.")
		http.Error(w, "Missing access token.", http.StatusBadRequest)
		return
	}
	readUser, err := wh.UserInteractor.Read(accessToken, wh.OauthConfiguration)
	if err != nil {
		wh.Logger.Log("[web_ReadUser] " + err.Error())
		http.Error(w, "Unnable to read user.", http.StatusInternalServerError)
		return
	}
	user, err := wh.UserInteractor.New(readUser.Name, readUser.Email, readUser.Login, accessToken)
	if err != nil {
		wh.Logger.Log("[web_ReadUser] " + err.Error())
		http.Error(w, "Unnable to store new user.", http.StatusInternalServerError)
		return
	}
	response, _ := json.Marshal(user)
	wh.Logger.Log("[web_ReadUser] Read user: " + string(response))
	w.WriteHeader(http.StatusOK)
	w.Write(response)

}

// ValidateRepository checks wether a given repository is a valid webpage project or not.
func (wh *WebHandler) ValidateRepository(w http.ResponseWriter, r *http.Request) {
	wh.Logger.Log("[web_ValidateRepository]")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	name := r.FormValue("name")
	projectBranch := r.FormValue("projectBranch")
	publicBranch := r.FormValue("publicBranch")
	accessToken := r.FormValue("accessToken")
	if accessToken == "" || name == "" || projectBranch == "" || publicBranch == "" {
		wh.Logger.Log("[web_ValidateRepository] Missing required parameters.")
		http.Error(w, "Missing required parameters.", http.StatusBadRequest)
		return
	}
	isValid, err := wh.RepositoryInteractor.Validate(name, projectBranch, accessToken, wh.OauthConfiguration)
	if err != nil {
		wh.Logger.Log("[web_ValidateRepository] " + err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !isValid {
		wh.Logger.Log("[web_ValidateRepository] Repository is not valid.")
		http.Error(w, "Repository is not valid.", http.StatusInternalServerError)
		return
	}
	_, err = wh.RepositoryInteractor.New(name, projectBranch, publicBranch, accessToken)
	if err != nil {
		wh.Logger.Log("[web_ValidateRepository] " + err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	wh.Logger.Log("[web_ValidateRepository] Repository is valid.")
	w.WriteHeader(http.StatusOK)
}

// CreateContent creates a new content file.
func (wh *WebHandler) CreateContent(w http.ResponseWriter, r *http.Request) {
	wh.Logger.Log("[web_CreateContent]")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	content, err := decodeJSON(r.Body, &usecases.Content{})
	if err != nil {
		wh.Logger.Log("[web_CreateContent] " + err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	createdContent, err := wh.ContentInteractor.New(*content.(*usecases.Content), wh.OauthConfiguration)
	if err != nil {
		wh.Logger.Log("[web_CreateContent] " + err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response, _ := json.Marshal(&createdContent)
	wh.Logger.Log("[web_CreateContent] Success with response: " + string(response))
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

// ReadContent reads the content of a file.
func (wh *WebHandler) ReadContent(w http.ResponseWriter, r *http.Request) {
	wh.Logger.Log("[web_ReadContent]")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	path := r.FormValue("path")
	projectBranch := r.FormValue("projectBranch")
	accessToken := r.FormValue("accessToken")
	repositoryName := r.FormValue("repositoryName")
	title := r.FormValue("title")
	if path == "" || projectBranch == "" || accessToken == "" || repositoryName == "" || title == "" {
		wh.Logger.Log("[web_ReadContent] Missing required parameters")
		http.Error(w, "Missing required parameters.", http.StatusBadRequest)
		return
	}
	content, err := wh.ContentInteractor.Find(usecases.Content{
		Path:           path,
		ProjectBranch:  projectBranch,
		RepositoryName: repositoryName,
		Title:          title,
		AccessToken:    accessToken,
	}, wh.OauthConfiguration)
	if err != nil {
		wh.Logger.Log("[web_ReadContent] " + err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response, _ := json.Marshal(content)
	wh.Logger.Log("[web_ReadContent] Success with response: " + string(response))
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

// SaveContent updates the content file.
func (wh *WebHandler) SaveContent(w http.ResponseWriter, r *http.Request) {
	wh.Logger.Log("[web_UpdateContent]")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	content, err := decodeJSON(r.Body, &usecases.Content{})
	if err != nil {
		wh.Logger.Log("[web_UpdateContent] " + err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	updatedContent, err := wh.ContentInteractor.Update(*content.(*usecases.Content), wh.OauthConfiguration)
	if err != nil {
		wh.Logger.Log("[web_UpdateContent] " + err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response, _ := json.Marshal(updatedContent)
	wh.Logger.Log("[web_UpdateContent] Success with response: " + string(response))
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

// ListContentTitles lists the content elements.
func (wh *WebHandler) ListContent(w http.ResponseWriter, r *http.Request) {
	wh.Logger.Log("[web_ListContent]")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	path := r.FormValue("path")
	projectBranch := r.FormValue("projectBranch")
	repositoryName := r.FormValue("repositoryName")
	accessToken := r.FormValue("accessToken")
	if path == "" || projectBranch == "" || accessToken == "" || repositoryName == "" {
		wh.Logger.Log("[web_ListContent] Missing required parameters")
		http.Error(w, "Missing required parameters.", http.StatusBadRequest)
		return
	}
	contents, err := wh.ContentInteractor.List(usecases.Content{
		Path:           path,
		ProjectBranch:  projectBranch,
		RepositoryName: repositoryName,
		AccessToken:    accessToken,
	}, wh.OauthConfiguration)
	if err != nil {
		wh.Logger.Log("[web_ListContent] " + err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response, _ := json.Marshal(contents)
	wh.Logger.Log("[web_ListContent] Success with response: " + string(response))
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

// RemoveContent removes a content file.
func (wh *WebHandler) RemoveContent(w http.ResponseWriter, r *http.Request) {
	wh.Logger.Log("[web_DeleteContent]")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	content, err := decodeJSON(r.Body, &usecases.Content{})
	if err != nil {
		wh.Logger.Log("[web_DeleteContent] " + err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = wh.ContentInteractor.Remove(*content.(*usecases.Content), wh.OauthConfiguration)
	if err != nil {
		wh.Logger.Log("[web_DeleteContent] " + err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	wh.Logger.Log("[web_DeleteContent] Sucess")
	w.WriteHeader(http.StatusOK)
}

// PublishContent publishes the content
func (wh *WebHandler) PublishContent(w http.ResponseWriter, r *http.Request) {
	wh.Logger.Log("[web_PublishContent]")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	content, err := decodeJSON(r.Body, &usecases.Content{})
	if err != nil {
		wh.Logger.Log("[web_PublishContent] " + err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = wh.ContentInteractor.Publish(*content.(*usecases.Content), wh.OauthConfiguration)
	if err != nil {
		wh.Logger.Log("[web_PublishContent] " + err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	wh.Logger.Log("[web_PublishContent] Success!")
	w.WriteHeader(http.StatusOK)
}

func (wh *WebHandler) isStateValid(state, receivedState string) bool {
	return strings.AreStringsEqual(state, receivedState)
}

func (wh *WebHandler) exchangeToken(code string, oauthConf *oauth2.Config) (*string, error) {
	token, err := oauthConf.Exchange(oauth2.NoContext, code)
	if err != nil {
		return nil, err
	}
	return &token.AccessToken, nil
}

func decodeJSON(body io.ReadCloser, structure interface{}) (interface{}, error) {
	if body == nil {
		return nil, errors.New("No request body found.")
	}
	err := json.NewDecoder(body).Decode(&structure)
	if err != nil {
		return nil, err
	}
	return structure, nil
}

package interfaces_test

import (
	"bytes"
	"encoding/json"
	"github.com/joaodias/hugito-backend/interfaces"
	"github.com/joaodias/hugito-backend/interfaces/mocks"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	mux                *http.ServeMux
	testServer         *httptest.Server
	oauthConfiguration *oauth2.Config
)

func setupExternalAuthenticationServer() {
	mux = http.NewServeMux()
	testServer = httptest.NewServer(mux)
	oauthConfiguration = &oauth2.Config{
		ClientID:     "CLIENT_ID",
		ClientSecret: "CLIENT_SECRET",
		Endpoint: oauth2.Endpoint{
			TokenURL: testServer.URL,
		},
	}
}

func TestAuthWebOptions(t *testing.T) {
	request, err := http.NewRequest("OPTIONS", "/auth", nil)
	if err != nil {
		t.Error("Error creating request.")
	}
	recorder := httptest.NewRecorder()
	webHandler := &interfaces.WebHandler{}
	webHandler.AuthWebOptions(recorder, request)
	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, "*", recorder.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "OPTIONS, POST", recorder.Header().Get("Access-Control-Allow-Methods"))
	assert.Equal(t, "application/json; charset=UTF-8", recorder.Header().Get("Content-Type"))
	assert.Equal(t, "access-control-allow-origin, content-type,  access-control-allow-headers, access-control-allow-methods", recorder.Header().Get("Access-Control-Allow-Headers"))
}

func TestContentWebOptions(t *testing.T) {
	request, err := http.NewRequest("OPTIONS", "/auth", nil)
	if err != nil {
		t.Error("Error creating request.")
	}
	recorder := httptest.NewRecorder()
	webHandler := &interfaces.WebHandler{}
	webHandler.ContentWebOptions(recorder, request)
	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, "*", recorder.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "OPTIONS, POST, GET, PUT, DELETE", recorder.Header().Get("Access-Control-Allow-Methods"))
	assert.Equal(t, "application/json; charset=UTF-8", recorder.Header().Get("Content-Type"))
	assert.Equal(t, "access-control-allow-origin, content-type,  access-control-allow-headers, access-control-allow-methods", recorder.Header().Get("Access-Control-Allow-Headers"))
}

func TestWebAuthenticateFailBadRequestNoBody(t *testing.T) {
	request, err := http.NewRequest("GET", "/auth", nil)
	if err != nil {
		t.Error("Error creating request.")
	}
	recorder := httptest.NewRecorder()
	webHandler := &interfaces.WebHandler{
		Logger: &mocks.Logger{},
	}
	webHandler.Authenticate(recorder, request)
	assert.Equal(t, recorder.Code, http.StatusBadRequest)
}

func TestWebAuthenticateFailMalformedJSON(t *testing.T) {
	body, err := json.Marshal(map[string]interface{}{
		"state": 1234,
	})
	request, err := http.NewRequest("POST", "/auth", bytes.NewReader(body))
	if err != nil {
		t.Error("Error creating request.")
	}
	recorder := httptest.NewRecorder()
	webHandler := &interfaces.WebHandler{
		Logger: &mocks.Logger{},
	}
	webHandler.Authenticate(recorder, request)
	assert.Equal(t, recorder.Code, http.StatusBadRequest)
}

func TestWebAuthenticationFailWhenDifferentStates(t *testing.T) {
	body, err := json.Marshal(map[string]interface{}{
		"state":         "state",
		"code":          "code",
		"receivedState": "receivedState",
	})
	request, err := http.NewRequest("POST", "/content", bytes.NewReader(body))
	if err != nil {
		t.Error("Error creating request.")
	}
	recorder := httptest.NewRecorder()
	webHandler := &interfaces.WebHandler{
		Logger: &mocks.Logger{},
	}
	webHandler.Authenticate(recorder, request)
	assert.Equal(t, recorder.Code, http.StatusNotAcceptable)
}

func TestWebAuthenticationFailWhenCantGetToken(t *testing.T) {
	setupExternalAuthenticationServer()
	body, err := json.Marshal(map[string]interface{}{
		"state":         "state",
		"code":          "code",
		"receivedState": "state",
	})
	request, err := http.NewRequest("POST", "/content", bytes.NewReader(body))
	if err != nil {
		t.Error("Error creating request.")
	}
	recorder := httptest.NewRecorder()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Unnable to exchange token.", http.StatusBadRequest)
	})
	webHandler := &interfaces.WebHandler{
		OauthConfiguration: oauthConfiguration,
		Logger:             &mocks.Logger{},
	}
	webHandler.Authenticate(recorder, request)
	assert.Equal(t, recorder.Code, http.StatusInternalServerError)
}

func TestWebAuthenticateSuccess(t *testing.T) {
	setupExternalAuthenticationServer()
	body, err := json.Marshal(map[string]interface{}{
		"state":         "state",
		"code":          "code",
		"receivedState": "state",
	})
	request, err := http.NewRequest("POST", "/content", bytes.NewReader(body))
	if err != nil {
		t.Error("Error creating request.")
	}
	recorder := httptest.NewRecorder()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("accessToken"))
		w.WriteHeader(http.StatusOK)
	})
	webHandler := &interfaces.WebHandler{
		OauthConfiguration: oauthConfiguration,
		Logger:             &mocks.Logger{},
	}
	webHandler.Authenticate(recorder, request)
	assert.Equal(t, recorder.Code, http.StatusOK)
}

func TestWebReadUserFailBadRequest(t *testing.T) {
	request, err := http.NewRequest("GET", "/user", nil)
	if err != nil {
		t.Error("Error creating request.")
	}
	recorder := httptest.NewRecorder()
	webHandler := &interfaces.WebHandler{
		Logger: &mocks.Logger{},
	}
	webHandler.ReadUser(recorder, request)
	assert.Equal(t, recorder.Code, http.StatusBadRequest)
}

func TestWebReadUserFailWhenReading(t *testing.T) {
	request, err := http.NewRequest("GET", "/user?accessToken=1234", nil)
	if err != nil {
		t.Error("Error creating request.")
	}
	recorder := httptest.NewRecorder()
	mockUserInteractor := &mocks.UserInteractor{
		IsReadError: true,
		IsNewError:  false,
	}
	webHandler := &interfaces.WebHandler{
		UserInteractor:     mockUserInteractor,
		OauthConfiguration: &oauth2.Config{},
		Logger:             &mocks.Logger{},
	}
	webHandler.ReadUser(recorder, request)
	assert.Equal(t, recorder.Code, http.StatusInternalServerError)
	assert.True(t, mockUserInteractor.IsReadCalled)
	assert.False(t, mockUserInteractor.IsNewCalled)
}

func TestWebReadUserFailWhenCreatingNewUser(t *testing.T) {
	request, err := http.NewRequest("GET", "/user?accessToken=1234", nil)
	if err != nil {
		t.Error("Error creating request.")
	}
	recorder := httptest.NewRecorder()
	mockUserInteractor := &mocks.UserInteractor{
		IsReadError: false,
		IsNewError:  true,
	}
	webHandler := &interfaces.WebHandler{
		UserInteractor:     mockUserInteractor,
		OauthConfiguration: &oauth2.Config{},
		Logger:             &mocks.Logger{},
	}
	webHandler.ReadUser(recorder, request)
	assert.Equal(t, recorder.Code, http.StatusInternalServerError)
	assert.True(t, mockUserInteractor.IsReadCalled)
	assert.True(t, mockUserInteractor.IsNewCalled)
}

func TestWebReadUserSuccess(t *testing.T) {
	request, err := http.NewRequest("GET", "/user?accessToken=1234", nil)
	if err != nil {
		t.Error("Error creating request.")
	}
	recorder := httptest.NewRecorder()
	mockUserInteractor := &mocks.UserInteractor{
		IsReadError: false,
		IsNewError:  false,
	}
	webHandler := &interfaces.WebHandler{
		UserInteractor:     mockUserInteractor,
		OauthConfiguration: &oauth2.Config{},
		Logger:             &mocks.Logger{},
	}
	webHandler.ReadUser(recorder, request)
	assert.Equal(t, recorder.Code, http.StatusOK)
	assert.True(t, mockUserInteractor.IsReadCalled)
	assert.True(t, mockUserInteractor.IsNewCalled)
}

func TestWebValidateRepositoryFailBadRequest(t *testing.T) {
	request, err := http.NewRequest("GET", "/repository", nil)
	if err != nil {
		t.Error("Error creating request.")
	}
	recorder := httptest.NewRecorder()
	webHandler := &interfaces.WebHandler{
		Logger: &mocks.Logger{},
	}
	webHandler.ValidateRepository(recorder, request)
	assert.Equal(t, recorder.Code, http.StatusBadRequest)
}

func TestWebValidateRepositoryFailWhenValidating(t *testing.T) {
	request, err := http.NewRequest("GET", "/repository?name=name&projectBranch=branch&publicBranch=pb&accessToken=at", nil)
	if err != nil {
		t.Error("Error creating request.")
	}
	recorder := httptest.NewRecorder()
	mockRepositoryInteractor := &mocks.RepositoryInteractor{
		IsValidateError: true,
		IsNewError:      false,
	}
	webHandler := &interfaces.WebHandler{
		RepositoryInteractor: mockRepositoryInteractor,
		Logger:               &mocks.Logger{},
	}
	webHandler.ValidateRepository(recorder, request)
	assert.True(t, mockRepositoryInteractor.IsValidateCalled)
	assert.False(t, mockRepositoryInteractor.IsNewCalled)
	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
}

func TestWebValidateRepositorySuccessWhenInvalidRepository(t *testing.T) {
	request, err := http.NewRequest("GET", "/repository?name=name&projectBranch=branch&publicBranch=pb&accessToken=at", nil)
	if err != nil {
		t.Error("Error creating request.")
	}
	recorder := httptest.NewRecorder()
	mockRepositoryInteractor := &mocks.RepositoryInteractor{
		IsValidRepository: false,
	}
	webHandler := &interfaces.WebHandler{
		RepositoryInteractor: mockRepositoryInteractor,
		Logger:               &mocks.Logger{},
	}
	webHandler.ValidateRepository(recorder, request)
	assert.True(t, mockRepositoryInteractor.IsValidateCalled)
	assert.False(t, mockRepositoryInteractor.IsNewCalled)
	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
}

func TestWebValidateRepositoryFailWhenStoringValidated(t *testing.T) {
	request, err := http.NewRequest("GET", "/repository?name=name&projectBranch=branch&publicBranch=pb&accessToken=at", nil)
	if err != nil {
		t.Error("Error creating request.")
	}
	recorder := httptest.NewRecorder()
	mockRepositoryInteractor := &mocks.RepositoryInteractor{
		IsValidateError:   false,
		IsNewError:        true,
		IsValidRepository: true,
	}
	webHandler := &interfaces.WebHandler{
		RepositoryInteractor: mockRepositoryInteractor,
		Logger:               &mocks.Logger{},
	}
	webHandler.ValidateRepository(recorder, request)
	assert.True(t, mockRepositoryInteractor.IsValidateCalled)
	assert.True(t, mockRepositoryInteractor.IsNewCalled)
	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
}

func TestWebValidateRepositorySuccessWhenValidRepository(t *testing.T) {
	request, err := http.NewRequest("GET", "/repository?name=name&projectBranch=branch&publicBranch=pb&accessToken=at", nil)
	if err != nil {
		t.Error("Error creating request.")
	}
	recorder := httptest.NewRecorder()
	mockRepositoryInteractor := &mocks.RepositoryInteractor{
		IsValidRepository: true,
	}
	webHandler := &interfaces.WebHandler{
		RepositoryInteractor: mockRepositoryInteractor,
		Logger:               &mocks.Logger{},
	}
	webHandler.ValidateRepository(recorder, request)
	assert.True(t, mockRepositoryInteractor.IsValidateCalled)
	assert.True(t, mockRepositoryInteractor.IsNewCalled)
	assert.Equal(t, http.StatusOK, recorder.Code)
}

func TestWebReadContentFailBadRequest(t *testing.T) {
	request, err := http.NewRequest("GET", "/content", nil)
	if err != nil {
		t.Error("Error creating request.")
	}
	recorder := httptest.NewRecorder()
	webHandler := &interfaces.WebHandler{
		Logger: &mocks.Logger{},
	}
	webHandler.ReadContent(recorder, request)
	assert.Equal(t, recorder.Code, http.StatusBadRequest)
}

func TestWebReadContentFailWhenReadingContent(t *testing.T) {
	request, err := http.NewRequest("GET", "/content?path=path&title=title&repositoryName=rn&projectBranch=branch&accessToken=accessToken", nil)
	if err != nil {
		t.Error("Error creating request.")
	}
	recorder := httptest.NewRecorder()
	mockContentInteractor := &mocks.ContentInteractor{
		IsFindError: true,
	}
	webHandler := &interfaces.WebHandler{
		ContentInteractor: mockContentInteractor,
		Logger:            &mocks.Logger{},
	}
	webHandler.ReadContent(recorder, request)
	assert.Equal(t, recorder.Code, http.StatusInternalServerError)
	assert.True(t, mockContentInteractor.FindCalled)
}

func TestWebReadContentSuccess(t *testing.T) {
	request, err := http.NewRequest("GET", "/content?path=path&title=title&repositoryName=rn&projectBranch=branch&accessToken=accessToken", nil)
	if err != nil {
		t.Error("Error creating request.")
	}
	recorder := httptest.NewRecorder()
	mockContentInteractor := &mocks.ContentInteractor{}
	webHandler := &interfaces.WebHandler{
		ContentInteractor: mockContentInteractor,
		Logger:            &mocks.Logger{},
	}
	webHandler.ReadContent(recorder, request)
	assert.Equal(t, recorder.Code, http.StatusOK)
	assert.True(t, mockContentInteractor.FindCalled)
}

func TestWebListContentTitlesFailBadRequest(t *testing.T) {
	request, err := http.NewRequest("GET", "/content", nil)
	if err != nil {
		t.Error("Error creating request.")
	}
	recorder := httptest.NewRecorder()
	webHandler := &interfaces.WebHandler{
		Logger: &mocks.Logger{},
	}
	webHandler.ListContent(recorder, request)
	assert.Equal(t, recorder.Code, http.StatusBadRequest)
}

func TestWebListContentTitlesFailWhenUnnableToListContent(t *testing.T) {
	request, err := http.NewRequest("GET", "/content?accessToken=at&projectBranch=pb&repositoryName=rn&path=path", nil)
	if err != nil {
		t.Error("Error creating request.")
	}
	recorder := httptest.NewRecorder()
	mockContentInteractor := &mocks.ContentInteractor{
		IsListError: true,
	}
	webHandler := &interfaces.WebHandler{
		OauthConfiguration: &oauth2.Config{},
		ContentInteractor:  mockContentInteractor,
		Logger:             &mocks.Logger{},
	}
	webHandler.ListContent(recorder, request)
	assert.True(t, mockContentInteractor.ListCalled)
	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
}

func TestWebListContentTitlesSuccess(t *testing.T) {
	request, err := http.NewRequest("GET", "/content?accessToken=at&projectBranch=pb&repositoryName=rn&path=path", nil)
	if err != nil {
		t.Error("Error creating request.")
	}
	recorder := httptest.NewRecorder()
	mockContentInteractor := &mocks.ContentInteractor{}
	webHandler := &interfaces.WebHandler{
		OauthConfiguration: &oauth2.Config{},
		ContentInteractor:  mockContentInteractor,
		Logger:             &mocks.Logger{},
	}
	webHandler.ListContent(recorder, request)
	assert.True(t, mockContentInteractor.ListCalled)
	assert.Equal(t, recorder.Code, http.StatusOK)
}

func TestWebCreateContentFailBadRequestNoBody(t *testing.T) {
	request, err := http.NewRequest("POST", "/content", nil)
	if err != nil {
		t.Error("Error creating request.")
	}
	recorder := httptest.NewRecorder()
	webHandler := &interfaces.WebHandler{
		Logger: &mocks.Logger{},
	}
	webHandler.CreateContent(recorder, request)
	assert.Equal(t, recorder.Code, http.StatusBadRequest)
}

func TestWebCreateContentFailMalformedJSON(t *testing.T) {
	// Repository name should be a string. This should trigger a json decoding error.
	body, err := json.Marshal(map[string]interface{}{
		"repositoryName": 1234,
	})
	request, err := http.NewRequest("POST", "/content", bytes.NewReader(body))
	if err != nil {
		t.Error("Error creating request.")
	}
	recorder := httptest.NewRecorder()
	webHandler := &interfaces.WebHandler{
		Logger: &mocks.Logger{},
	}
	webHandler.CreateContent(recorder, request)
	assert.Equal(t, recorder.Code, http.StatusBadRequest)
}

func TestWebCreateContentFailCreatingContent(t *testing.T) {
	body, err := json.Marshal(map[string]interface{}{
		"repositoryName": "repositoryName",
	})
	request, err := http.NewRequest("POST", "/content", bytes.NewReader(body))
	if err != nil {
		t.Error("Error creating request.")
	}
	recorder := httptest.NewRecorder()
	mockContentInteractor := &mocks.ContentInteractor{
		IsNewError: true,
	}
	webHandler := &interfaces.WebHandler{
		ContentInteractor: mockContentInteractor,
		Logger:            &mocks.Logger{},
	}
	webHandler.CreateContent(recorder, request)
	assert.True(t, mockContentInteractor.NewCalled)
	assert.Equal(t, recorder.Code, http.StatusInternalServerError)
}

func TestWebCreateContentSuccess(t *testing.T) {
	body, err := json.Marshal(map[string]interface{}{
		"repositoryName": "repositoryName",
		"projectBranch":  "projectBranch",
		"title":          "title",
		"path":           "path",
		"body":           "body",
	})
	request, err := http.NewRequest("POST", "/content", bytes.NewReader(body))
	if err != nil {
		t.Error("Error creating request.")
	}
	recorder := httptest.NewRecorder()
	mockContentInteractor := &mocks.ContentInteractor{}
	webHandler := &interfaces.WebHandler{
		ContentInteractor: mockContentInteractor,
		Logger:            &mocks.Logger{},
	}
	webHandler.CreateContent(recorder, request)
	assert.True(t, mockContentInteractor.NewCalled)
	assert.Equal(t, http.StatusOK, recorder.Code)
}

func TestWebSaveContentFailBadRequestNoBody(t *testing.T) {
	request, err := http.NewRequest("PUT", "/content", nil)
	if err != nil {
		t.Error("Error creating request.")
	}
	recorder := httptest.NewRecorder()
	webHandler := &interfaces.WebHandler{
		Logger: &mocks.Logger{},
	}
	webHandler.SaveContent(recorder, request)
	assert.Equal(t, recorder.Code, http.StatusBadRequest)
}

func TestWebSaveContentFailMalformedJSON(t *testing.T) {
	body, err := json.Marshal(map[string]interface{}{
		"repositoryName": 1234,
	})
	request, err := http.NewRequest("PUT", "/content", bytes.NewReader(body))
	if err != nil {
		t.Error("Error creating request.")
	}
	recorder := httptest.NewRecorder()
	webHandler := &interfaces.WebHandler{
		Logger: &mocks.Logger{},
	}
	webHandler.SaveContent(recorder, request)
	assert.Equal(t, recorder.Code, http.StatusBadRequest)
}

func TestWebSaveContentFailWhenUpdatingContent(t *testing.T) {
	body, err := json.Marshal(map[string]interface{}{
		"repositoryName": "repositoryName",
		"projectBranch":  "projectBranch",
		"title":          "title",
		"path":           "path",
		"body":           "body",
	})
	request, err := http.NewRequest("PUT", "/content", bytes.NewReader(body))
	if err != nil {
		t.Error("Error creating request.")
	}
	recorder := httptest.NewRecorder()
	mockContentInteractor := &mocks.ContentInteractor{
		IsUpdateError: true,
	}
	webHandler := &interfaces.WebHandler{
		ContentInteractor: mockContentInteractor,
		Logger:            &mocks.Logger{},
	}
	webHandler.SaveContent(recorder, request)
	assert.Equal(t, recorder.Code, http.StatusInternalServerError)
}

func TestWebSaveContentSuccess(t *testing.T) {
	body, err := json.Marshal(map[string]interface{}{
		"repositoryName": "repositoryName",
		"projectBranch":  "projectBranch",
		"title":          "title",
		"path":           "path",
		"body":           "body",
	})
	request, err := http.NewRequest("POST", "/content", bytes.NewReader(body))
	if err != nil {
		t.Error("Error creating request.")
	}
	recorder := httptest.NewRecorder()
	mockContentInteractor := &mocks.ContentInteractor{}
	webHandler := &interfaces.WebHandler{
		ContentInteractor: mockContentInteractor,
		Logger:            &mocks.Logger{},
	}
	webHandler.SaveContent(recorder, request)
	assert.Equal(t, recorder.Code, http.StatusOK)
}

func TestWebRemoveContentFailBadRequestNoBody(t *testing.T) {
	request, err := http.NewRequest("DELETE", "/content", nil)
	if err != nil {
		t.Error("Error creating request.")
	}
	recorder := httptest.NewRecorder()
	webHandler := &interfaces.WebHandler{
		Logger: &mocks.Logger{},
	}
	webHandler.RemoveContent(recorder, request)
	assert.Equal(t, recorder.Code, http.StatusBadRequest)
}

func TestWebRemoveContentFailMalformedJSON(t *testing.T) {
	body, err := json.Marshal(map[string]interface{}{
		"repositoryName": 1234,
	})
	request, err := http.NewRequest("PUT", "/content", bytes.NewReader(body))
	if err != nil {
		t.Error("Error creating request.")
	}
	recorder := httptest.NewRecorder()
	webHandler := &interfaces.WebHandler{
		Logger: &mocks.Logger{},
	}
	webHandler.RemoveContent(recorder, request)
	assert.Equal(t, recorder.Code, http.StatusBadRequest)
}

func TestWebRemoveContentFailRemovingContent(t *testing.T) {
	body, err := json.Marshal(map[string]interface{}{
		"id": "id",
	})
	request, err := http.NewRequest("DELETE", "/content", bytes.NewReader(body))
	if err != nil {
		t.Error("Error creating request.")
	}
	recorder := httptest.NewRecorder()
	mockContentInteractor := &mocks.ContentInteractor{
		IsRemoveError: true,
	}
	webHandler := &interfaces.WebHandler{
		ContentInteractor: mockContentInteractor,
		Logger:            &mocks.Logger{},
	}
	webHandler.RemoveContent(recorder, request)
	assert.True(t, mockContentInteractor.RemoveCalled)
	assert.Equal(t, recorder.Code, http.StatusInternalServerError)
}

func TestWebRemoveContentSuccess(t *testing.T) {
	body, err := json.Marshal(map[string]interface{}{
		"id": "id",
	})
	request, err := http.NewRequest("DELETE", "/content", bytes.NewReader(body))
	if err != nil {
		t.Error("Error creating request.")
	}
	recorder := httptest.NewRecorder()
	mockContentInteractor := &mocks.ContentInteractor{}
	webHandler := &interfaces.WebHandler{
		ContentInteractor: mockContentInteractor,
		Logger:            &mocks.Logger{},
	}
	webHandler.RemoveContent(recorder, request)
	assert.True(t, mockContentInteractor.RemoveCalled)
	assert.Equal(t, recorder.Code, http.StatusOK)
}

func TestWebPublishContentFailBadRequestNoBody(t *testing.T) {
	request, err := http.NewRequest("POST", "/content/publish", nil)
	if err != nil {
		t.Error("Error creating request.")
	}
	recorder := httptest.NewRecorder()
	webHandler := &interfaces.WebHandler{
		Logger: &mocks.Logger{},
	}
	webHandler.PublishContent(recorder, request)
	assert.Equal(t, recorder.Code, http.StatusBadRequest)
}

func TestWebPublishContentFailMalformedJSON(t *testing.T) {
	body, err := json.Marshal(map[string]interface{}{
		"repositoryName": 1234,
	})
	request, err := http.NewRequest("PUT", "/content/publish", bytes.NewReader(body))
	if err != nil {
		t.Error("Error creating request.")
	}
	recorder := httptest.NewRecorder()
	webHandler := &interfaces.WebHandler{
		Logger: &mocks.Logger{},
	}
	webHandler.PublishContent(recorder, request)
	assert.Equal(t, recorder.Code, http.StatusBadRequest)
}

func TestWebPublishContentFailPublishingContent(t *testing.T) {
	body, err := json.Marshal(map[string]interface{}{
		"id": "id",
	})
	request, err := http.NewRequest("POST", "/content/publish", bytes.NewReader(body))
	if err != nil {
		t.Error("Error creating request.")
	}
	recorder := httptest.NewRecorder()
	mockContentInteractor := &mocks.ContentInteractor{
		IsPublishError: true,
	}
	webHandler := &interfaces.WebHandler{
		ContentInteractor: mockContentInteractor,
		Logger:            &mocks.Logger{},
	}
	webHandler.PublishContent(recorder, request)
	assert.True(t, mockContentInteractor.PublishCalled)
	assert.Equal(t, recorder.Code, http.StatusInternalServerError)
}

func TestWebPublishContentSuccess(t *testing.T) {
	body, err := json.Marshal(map[string]interface{}{
		"id": "id",
	})
	request, err := http.NewRequest("POST", "/content/publish", bytes.NewReader(body))
	if err != nil {
		t.Error("Error creating request.")
	}
	recorder := httptest.NewRecorder()
	mockContentInteractor := &mocks.ContentInteractor{}
	webHandler := &interfaces.WebHandler{
		ContentInteractor: mockContentInteractor,
		Logger:            &mocks.Logger{},
	}
	webHandler.PublishContent(recorder, request)
	assert.True(t, mockContentInteractor.PublishCalled)
	assert.Equal(t, recorder.Code, http.StatusOK)
}

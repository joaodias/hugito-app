package handlers_test

import (
	"fmt"
	"github.com/fatih/structs"
	handlers "github.com/joaodias/hugito-app/handlers"
	mocks "github.com/joaodias/hugito-app/mocks"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"golang.org/x/oauth2"
	"net/http"
	"net/http/httptest"
)

var _ = Describe("Handlers", func() {
	const (
		UserFinished = iota
		AuthenticationFinished
		RepositoryFinished
		ContentFinished
	)

	Describe("Authentication Handlers", func() {
		type MockAuthentication struct {
			Authenticated string `json:"authenticated"`
			Code          string `json:"code"`
			State         string `json:"state"`
			ReceivedState string `json:"receivedState"`
		}
		var (
			mClient    *mocks.Client
			mux        *http.ServeMux
			testServer *httptest.Server
		)
		BeforeEach(func() {
			mux = http.NewServeMux()
			testServer = httptest.NewServer(mux)
			mClient = &mocks.Client{
				FinishedChannels: make(map[int]chan bool),
				OauthConf: &oauth2.Config{
					ClientID:     "FAKE_CLIENT_ID",
					ClientSecret: "FAKE_CLIENT_SECRET",
					Endpoint: oauth2.Endpoint{
						TokenURL: testServer.URL,
					},
				},
			}
		})
		Describe("When authentication is requested", func() {
			Context("and the JSON is invalid", func() {
				It("should send to the client an error message", func() {
					handlers.Authenticate(mClient, "some stuff that looks like an invalid json")
					Expect(mClient.Name).To(ContainSubstring("error"))
					Expect(mClient.Data).To(ContainSubstring("Error decoding json:"))
				})
			})
			Context("and the states are different", func() {
				It("should return an error", func() {
					mockAuthentication := &MockAuthentication{Authenticated: "false", Code: "1234", State: "5678", ReceivedState: "9867"}
					mData := structs.Map(mockAuthentication)
					handlers.Authenticate(mClient, mData)
					<-mClient.FinishedChannels[AuthenticationFinished]
					Expect(mClient.Name).To(ContainSubstring("error"))
					Expect(mClient.Data).To(ContainSubstring("received state and state are different."))
				})
			})
			Context("and couldn't get token", func() {
				It("should send to the client an error message", func() {
					defer testServer.Close()
					mockAuthentication := &MockAuthentication{Authenticated: "false", Code: "1234", State: "5678", ReceivedState: "5678"}
					mData := structs.Map(mockAuthentication)
					// the server is handlers is not defined so
					// it will return a not found, which is cool.
					handlers.Authenticate(mClient, mData)
					<-mClient.FinishedChannels[AuthenticationFinished]
					Expect(mClient.Name).To(ContainSubstring("error"))
					Expect(mClient.Data).To(ContainSubstring("Error getting the access token."))
				})
			})
			Context("and successfully got the token", func() {
				It("should send to the client a authenticated set message with the token", func() {
					mockAuthentication := &MockAuthentication{Authenticated: "false", Code: "1234", State: "5678", ReceivedState: "5678"}
					mData := structs.Map(mockAuthentication)
					defer testServer.Close()
					mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
						fmt.Fprint(w, "access_token=90d64460d14870c08c81352a05dedd3465940a7c&scope=user&token_type=bearer")
					})
					handlers.Authenticate(mClient, mData)
					<-mClient.FinishedChannels[AuthenticationFinished]
					Expect(mClient.Name).To(ContainSubstring("authenticated set"))
					Expect(mClient.Data).To(ContainSubstring("90d64460d14870c08c81352a05dedd3465940a7c"))
				})
			})
		})
	})

	Describe("User Handlers", func() {
		type MockUser struct {
			Name        string `json:"name"`
			AccessToken string `json:"accessToken"`
		}
		var (
			mClient                *mocks.Client
			mux                    *http.ServeMux
			testServer             *httptest.Server
			mockWithValidToken     *MockUser
			mockJSONWithValidToken map[string]interface{}
		)
		BeforeEach(func() {
			mux = http.NewServeMux()
			testServer = httptest.NewServer(mux)
			mClient = &mocks.Client{
				FinishedChannels: make(map[int]chan bool),
				OauthConf: &oauth2.Config{
					ClientID:     "FAKE_CLIENT_ID",
					ClientSecret: "FAKE_CLIENT_SECRET",
					Endpoint: oauth2.Endpoint{
						TokenURL: testServer.URL,
					},
				},
			}
			mockWithValidToken = &MockUser{Name: "joaodias", AccessToken: "90d64460d14870c08c81352a05dedd3465940a7c"}
			mockJSONWithValidToken = structs.Map(mockWithValidToken)
		})
		Describe("When a get user is requested", func() {
			Context("and the json is invalid", func() {
				It("should send to the client an error message", func() {
					handlers.GetUser(mClient, "some stuff that looks like an invalid json")
					Expect(mClient.Name).To(ContainSubstring("error"))
					Expect(mClient.Data).To(ContainSubstring("Error decoding json:"))
				})
			})
			Context("and the user is not retrieved", func() {
				It("should send a logout message to the client", func() {
					defer testServer.Close()
					mux.HandleFunc("/user/", func(w http.ResponseWriter, r *http.Request) {
						fmt.Fprint(w, `somethingReallyWrong`)
					})
					handlers.GetUser(mClient, mockJSONWithValidToken)
					<-mClient.FinishedChannels[UserFinished]
					Expect(mClient.Name).To(ContainSubstring("logout"))
					Expect(mClient.Data).To(Equal("Cannot get the authorized user."))
				})
			})
			Context("and the user is successfully retrieved", func() {
				It("should return a user set message to the client", func() {
					defer testServer.Close()
					mux.HandleFunc("/user/", func(w http.ResponseWriter, r *http.Request) {
						fmt.Fprint(w, `{"Name":"joaodias"}`)
					})
					handlers.GetUser(mClient, mockJSONWithValidToken)
					<-mClient.FinishedChannels[UserFinished]
					receivedData := mClient.Data.(handlers.User)
					Expect(mClient.Name).To(ContainSubstring("user set"))
					Expect(receivedData.Name).To(Equal("joaodias"))
					Expect(receivedData.AccessToken).To(Equal("90d64460d14870c08c81352a05dedd3465940a7c"))
				})
			})
		})
	})

	Describe("Repository Handlers", func() {
		type MockRepositories struct {
			Names       []string `json:"names"`
			AccessToken string   `json:"accessToken"`
		}
		type MockRepository struct {
			Name        string `json:"name"`
			AccessToken string `json:"accessToken"`
		}
		var (
			mClient                *mocks.Client
			mux                    *http.ServeMux
			testServer             *httptest.Server
			mockWithValidToken     *MockRepositories
			mockJSONWithValidToken map[string]interface{}
			mockValidRepo          *MockRepository
			mockJSONValidRepo      map[string]interface{}
			mockInvalidRepo        *MockRepository
			mockJSONInvalidRepo    map[string]interface{}
		)
		BeforeEach(func() {
			mux = http.NewServeMux()
			testServer = httptest.NewServer(mux)
			mClient = &mocks.Client{
				FinishedChannels: make(map[int]chan bool),
				OauthConf: &oauth2.Config{
					ClientID:     "FAKE_CLIENT_ID",
					ClientSecret: "FAKE_CLIENT_SECRET",
					Endpoint: oauth2.Endpoint{
						TokenURL: testServer.URL,
					},
				},
			}
			mockWithValidToken = &MockRepositories{Names: []string{""},
				AccessToken: "90d64460d14870c08c81352a05dedd3465940a7c"}
			mockJSONWithValidToken = structs.Map(mockWithValidToken)
			mockValidRepo = &MockRepository{Name: "validrepo", AccessToken: "90d64460d14870c08c81352a05dedd3465940a7c"}
			mockJSONValidRepo = structs.Map(mockValidRepo)
			mockInvalidRepo = &MockRepository{Name: "invalidrepo", AccessToken: "90d64460d14870c08c81352a05dedd3465940a7c"}
			mockJSONInvalidRepo = structs.Map(mockInvalidRepo)
		})
		Describe("When asking for user repositories", func() {
			Context("and the JSON is invalid", func() {
				It("should return an error to the client.", func() {
					handlers.GetRepository(mClient, "some stuff that looks like an invalid json")
					Expect(mClient.Name).To(ContainSubstring("error"))
					Expect(mClient.Data).To(ContainSubstring("Error decoding json:"))
				})
			})
			Context("and repositories are not retrieved", func() {
				It("should return an error to the client.", func() {
					defer testServer.Close()
					mux.HandleFunc("/user/repos/", func(w http.ResponseWriter, r *http.Request) {
						fmt.Fprint(w, `reallywrongresponse`)
					})
					handlers.GetRepository(mClient, mockJSONWithValidToken)
					<-mClient.FinishedChannels[RepositoryFinished]
					Expect(mClient.Name).To(ContainSubstring("logout"))
					Expect(mClient.Data).To(Equal("Cannot get the user repositories."))
				})
			})
			Context("and repositories are successfully retrieved", func() {
				It("should return the retrieved repositories to the client.", func() {
					mux.HandleFunc("/user/repos/", func(w http.ResponseWriter, r *http.Request) {
						fmt.Fprint(w, `[{"Name":"repo1"}, {"Name":"repo2"}, {"Name":"repo3"}]`)
					})
					handlers.GetRepository(mClient, mockJSONWithValidToken)
					<-mClient.FinishedChannels[RepositoryFinished]
					receivedData := mClient.Data.(handlers.Repositories)
					Expect(mClient.Name).To(ContainSubstring("repository set"))
					Expect(receivedData.Names).To(Equal([]string{"repo1", "repo2", "repo3"}))
					Expect(receivedData.AccessToken).To(Equal("90d64460d14870c08c81352a05dedd3465940a7c"))
				})
			})
		})

		Describe("When checking the validity of a repository", func() {
			Context("and the JSON is invalid", func() {
				It("should return an error to the client", func() {
					handlers.ValidateRepository(mClient, "some stuff that looks like an invalid json")
					Expect(mClient.Name).To(ContainSubstring("error"))
					Expect(mClient.Data).To(ContainSubstring("Error decoding json:"))
				})
			})
			Context("and the user Login cannot be retrieved", func() {
				It("should return a  logout message to the client", func() {
					defer testServer.Close()
					mux.HandleFunc("/user/", func(w http.ResponseWriter, r *http.Request) {
						fmt.Fprint(w, `someErronicThingHappened`)
					})
					handlers.ValidateRepository(mClient, mockJSONValidRepo)
					<-mClient.FinishedChannels[RepositoryFinished]
					Expect(mClient.Name).To(Equal("logout"))
					Expect(mClient.Data).To(Equal("Can't retrieve the authenticated user."))
				})
			})
			Context("and the repository tree cannot be retrieved", func() {
				It("should return an error message to the client", func() {
					defer testServer.Close()
					mux.HandleFunc("/user/", func(w http.ResponseWriter, r *http.Request) {
						fmt.Fprint(w, `{"Login":"joaodias"}`)
					})
					mux.HandleFunc("/repos/joaodias/validrepo/contents/", func(w http.ResponseWriter, r *http.Request) {
						fmt.Fprint(w, `someErroniousStuff`)
					})
					handlers.ValidateRepository(mClient, mockJSONValidRepo)
					<-mClient.FinishedChannels[RepositoryFinished]
					Expect(mClient.Name).To(Equal("error"))
					Expect(mClient.Data).To(Equal("Can't retrieve selected repository."))
				})
			})
			Context("and the repository is invalid", func() {
				It("should return error with repository not valid to the client", func() {
					defer testServer.Close()
					mux.HandleFunc("/user/", func(w http.ResponseWriter, r *http.Request) {
						fmt.Fprint(w, `{"Login":"joaodias"}`)
					})
					mux.HandleFunc("/repos/joaodias/invalidrepo/contents/", func(w http.ResponseWriter, r *http.Request) {
						fmt.Fprint(w, `[{"Name":"dir1"}, {"Name":"file1"}, {"Name":"dir2"}]`)
					})
					handlers.ValidateRepository(mClient, mockJSONInvalidRepo)
					<-mClient.FinishedChannels[RepositoryFinished]
					Expect(mClient.Name).To(Equal("error"))
					Expect(mClient.Data).To(Equal("Repository is not valid."))
				})
			})
			Context("and the repository is valid", func() {
				It("should return true to the client", func() {
					defer testServer.Close()
					mux.HandleFunc("/user/", func(w http.ResponseWriter, r *http.Request) {
						fmt.Fprint(w, `{"Login":"joaodias"}`)
					})
					mux.HandleFunc("/repos/joaodias/validrepo/contents/", func(w http.ResponseWriter, r *http.Request) {
						fmt.Fprint(w, `[{"Name":"content"}, {"Name":"config.toml"}, {"Name":"public"}, {"Name":"themes"}]`)
					})
					handlers.ValidateRepository(mClient, mockJSONValidRepo)
					<-mClient.FinishedChannels[RepositoryFinished]
					Expect(mClient.Name).To(Equal("repository validate"))
					Expect(mClient.Data).To(Equal("Repository is valid."))
				})
			})
		})
	})
	Describe("Content handlers", func() {
		type MockContentList struct {
			Name        string   `json:"name"`
			Titles      []string `json:"title"`
			AccessToken string   `json:"accessToken"`
		}
		var (
			mClient                   *mocks.Client
			mux                       *http.ServeMux
			testServer                *httptest.Server
			mockValidatedRepo         *MockContentList
			mockJSONWithValidatedRepo map[string]interface{}
		)
		BeforeEach(func() {
			mux = http.NewServeMux()
			testServer = httptest.NewServer(mux)
			mClient = &mocks.Client{
				FinishedChannels: make(map[int]chan bool),
				OauthConf: &oauth2.Config{
					ClientID:     "FAKE_CLIENT_ID",
					ClientSecret: "FAKE_CLIENT_SECRET",
					Endpoint: oauth2.Endpoint{
						TokenURL: testServer.URL,
					},
				},
			}
			mockValidatedRepo = &MockContentList{Name: "validatedrepo", Titles: []string{""}, AccessToken: "90d64460d14870c08c81352a05dedd3465940a7c"}
			mockJSONWithValidatedRepo = structs.Map(mockValidatedRepo)
		})
		Describe("When getting a list of a content files", func() {
			Context("and the JSON is invalid", func() {
				It("should return an error to the client", func() {
					handlers.GetContentList(mClient, "some stuff that looks like an invalid json")
					Expect(mClient.Name).To(ContainSubstring("error"))
					Expect(mClient.Data).To(ContainSubstring("Error decoding json:"))
				})
			})
			Context("and the user Login cannot be retrieved", func() {
				It("should return a logout message to the client", func() {
					defer testServer.Close()
					mux.HandleFunc("/user/", func(w http.ResponseWriter, r *http.Request) {
						fmt.Fprint(w, `someErronicThingHappened`)
					})
					handlers.GetContentList(mClient, mockJSONWithValidatedRepo)
					<-mClient.FinishedChannels[ContentFinished]
					Expect(mClient.Name).To(Equal("logout"))
					Expect(mClient.Data).To(Equal("Can't retrieve the authenticated user."))
				})
			})
			Context("and the content list cannot be retrieved", func() {
				It("should return an error message to the client", func() {
					defer testServer.Close()
					mux.HandleFunc("/user/", func(w http.ResponseWriter, r *http.Request) {
						fmt.Fprint(w, `{"Login":"joaodias"}`)
					})
					mux.HandleFunc("/repos/joaodias/validatedrepo/contents/content", func(w http.ResponseWriter, r *http.Request) {
						fmt.Fprint(w, `someErroniousStuff`)
					})
					handlers.GetContentList(mClient, mockJSONWithValidatedRepo)
					<-mClient.FinishedChannels[ContentFinished]
					Expect(mClient.Name).To(Equal("error"))
					Expect(mClient.Data).To(Equal("Can't retrieve the content list."))
				})
			})
			Context("and the content list is successfully retrieved", func() {
				It("should return a content list message to the client", func() {
					defer testServer.Close()
					mux.HandleFunc("/user/", func(w http.ResponseWriter, r *http.Request) {
						fmt.Fprint(w, `{"Login":"joaodias"}`)
					})
					mux.HandleFunc("/repos/joaodias/validatedrepo/contents/content", func(w http.ResponseWriter, r *http.Request) {
						fmt.Fprint(w, `[{"Name":"Content File 1"}, {"Name":"Content File 2"}]`)
					})
					handlers.GetContentList(mClient, mockJSONWithValidatedRepo)
					<-mClient.FinishedChannels[ContentFinished]
					receivedData := mClient.Data.(handlers.ContentList)
					Expect(receivedData.Name).To(Equal("validatedrepo"))
					Expect(receivedData.Titles).To(Equal([]string{"Content File 1", "Content File 2"}))
					Expect(receivedData.AccessToken).To(Equal("90d64460d14870c08c81352a05dedd3465940a7c"))
				})
			})
		})
	})
})

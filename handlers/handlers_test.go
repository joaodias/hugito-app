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
	)

	Describe("Authentication Handlers", func() {
		var (
			mClient    *mocks.Client
			mux        *http.ServeMux
			testServer *httptest.Server
		)
		type MockAuthentication struct {
			Authenticated string `json:"authenticated"`
			Code          string `json:"code"`
			State         string `json:"state"`
			ReceivedState string `json:"receivedState"`
		}
		BeforeEach(func() {
			mClient = &mocks.Client{
				FinishedChannels: make(map[int]chan bool),
			}
		})
		Describe("When authentication is requested", func() {
			Context("and the JSON is invalid", func() {
				expectedName := "error"
				expectedData := "Error decoding json: "
				It("should send to the client an error message", func() {
					handlers.Authenticate(mClient, "some stuff that looks like an invalid json")
					<-mClient.FinishedChannels[AuthenticationFinished]
					Expect(mClient.Name).To(ContainSubstring(expectedName))
					Expect(mClient.Data).To(ContainSubstring(expectedData))
				})
			})
			Context("and the states are different", func() {
				expectedName := "error"
				expectedData := "received state and state are different."
				It("should return an error", func() {
					mockAuthentication := &MockAuthentication{Authenticated: "false", Code: "1234", State: "5678", ReceivedState: "9867"}
					mData := structs.Map(mockAuthentication)
					handlers.Authenticate(mClient, mData)
					<-mClient.FinishedChannels[AuthenticationFinished]
					Expect(mClient.Name).To(ContainSubstring(expectedName))
					Expect(mClient.Data).To(ContainSubstring(expectedData))
				})
			})
			Context("and successfully got the token", func() {
				expectedName := "authenticated set"
				expectedData := "90d64460d14870c08c81352a05dedd3465940a7c"
				It("should send to the client a authenticated set message with the token", func() {
					mux = http.NewServeMux()
					testServer = httptest.NewServer(mux)
					defer testServer.Close()
					mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
						fmt.Fprint(w, "access_token=90d64460d14870c08c81352a05dedd3465940a7c&scope=user&token_type=bearer")
					})
					mClient.OauthConf = &oauth2.Config{
						ClientID:     "FAKE_CLIENT_ID",
						ClientSecret: "FAKE_CLIENT_SECRET",
						Endpoint: oauth2.Endpoint{
							TokenURL: testServer.URL + "/token",
						},
					}
					mockAuthentication := &MockAuthentication{Authenticated: "false", Code: "1234", State: "5678", ReceivedState: "5678"}
					mData := structs.Map(mockAuthentication)
					handlers.Authenticate(mClient, mData)
					<-mClient.FinishedChannels[AuthenticationFinished]
					Expect(mClient.Name).To(ContainSubstring(expectedName))
					Expect(mClient.Data).To(ContainSubstring(expectedData))
				})
			})
			Context("and couldn't get token", func() {
				expectedName := "error"
				expectedData := "Error getting the access token."
				It("should send to the client an error message", func() {
					mClient.OauthConf = &oauth2.Config{
						ClientID:     "FAKE_CLIENT_ID",
						ClientSecret: "FAKE_CLIENT_SECRET",
						Endpoint: oauth2.Endpoint{
							TokenURL: "something" + "/token",
						},
					}
					mockAuthentication := &MockAuthentication{Authenticated: "false", Code: "1234", State: "5678", ReceivedState: "5678"}
					mData := structs.Map(mockAuthentication)
					handlers.Authenticate(mClient, mData)
					<-mClient.FinishedChannels[AuthenticationFinished]
					Expect(mClient.Name).To(ContainSubstring(expectedName))
					Expect(mClient.Data).To(ContainSubstring(expectedData))
				})
			})
		})
	})

	Describe("User Handlers", func() {
		var (
			mClient    *mocks.Client
			mux        *http.ServeMux
			testServer *httptest.Server
		)
		type MockUser struct {
			Name        string `json:"name"`
			AccessToken string `json:"accessToken"`
		}
		BeforeEach(func() {
			mClient = &mocks.Client{
				FinishedChannels: make(map[int]chan bool),
			}
		})
		Describe("When a get user is requested", func() {
			Context("and the json is invalid", func() {
				expectedName := "error"
				expectedData := "Error decoding json: "
				It("should send to the client an error message", func() {
					handlers.GetUser(mClient, "some stuff that looks like an invalid json")
					<-mClient.FinishedChannels[UserFinished]
					Expect(mClient.Name).To(ContainSubstring(expectedName))
					Expect(mClient.Data).To(ContainSubstring(expectedData))
				})
			})
			Context("and the user is successfully retrieved", func() {
				expectedName := "user set"
				expectedData := &MockUser{Name: "joaodias", AccessToken: "90d64460d14870c08c81352a05dedd3465940a7c"}
				It("should return a user set message to the client", func() {
					mockWithValidToken := &MockUser{Name: "joaodias", AccessToken: "90d64460d14870c08c81352a05dedd3465940a7c"}
					mockJsonWithValidToken := structs.Map(mockWithValidToken)
					mux = http.NewServeMux()
					testServer = httptest.NewServer(mux)
					defer testServer.Close()
					mux.HandleFunc("/user/", func(w http.ResponseWriter, r *http.Request) {
						fmt.Fprint(w, `{"Name":"joaodias"}`)
					})
					mClient.OauthConf = &oauth2.Config{
						ClientID:     "FAKE_CLIENT_ID",
						ClientSecret: "FAKE_CLIENT_SECRET",
						Endpoint: oauth2.Endpoint{
							TokenURL: testServer.URL,
						},
					}
					handlers.GetUser(mClient, mockJsonWithValidToken)
					<-mClient.FinishedChannels[UserFinished]
					receivedData := mClient.Data.(handlers.User)
					Expect(mClient.Name).To(ContainSubstring(expectedName))
					Expect(receivedData.Name).To(Equal(expectedData.Name))
					Expect(receivedData.AccessToken).To(Equal(expectedData.AccessToken))
				})
			})
			Context("and the user is not retrieved", func() {
				expectedName := "logout"
				expectedData := "Cannot get the authorized user."
				It("should return a user set message to the client", func() {
					mockWithValidToken := &MockUser{Name: "joaodias", AccessToken: "90d64460d14870c08c81352a05dedd3465940a7c"}
					mockJsonWithValidToken := structs.Map(mockWithValidToken)
					mux = http.NewServeMux()
					testServer = httptest.NewServer(mux)
					defer testServer.Close()
					mux.HandleFunc("/user/", func(w http.ResponseWriter, r *http.Request) {
						fmt.Fprint(w, `{"Name":"joaodias"}`)
					})
					// To simulate the failed test a blank url is given to the
					// github client server. This way, the server can't get any // username and instead returns an error.
					mClient.OauthConf = &oauth2.Config{
						ClientID:     "FAKE_CLIENT_ID",
						ClientSecret: "FAKE_CLIENT_SECRET",
						Endpoint: oauth2.Endpoint{
							TokenURL: "",
						},
					}
					handlers.GetUser(mClient, mockJsonWithValidToken)
					<-mClient.FinishedChannels[UserFinished]
					Expect(mClient.Name).To(ContainSubstring(expectedName))
					Expect(mClient.Data).To(Equal(expectedData))
				})
			})
		})
	})

	Describe("Repository Handlers", func() {
		var (
			mClient    *mocks.Client
			mux        *http.ServeMux
			testServer *httptest.Server
		)
		type MockRepositories struct {
			Names       []string `json:"names"`
			AccessToken string   `json:"accessToken"`
		}
		BeforeEach(func() {
			mClient = &mocks.Client{
				FinishedChannels: make(map[int]chan bool),
			}
		})
		Describe("When asking for user repositories", func() {
			Context("and the JSON is invalid", func() {
				expectedName := "error"
				expectedData := "Error decoding json:"
				It("should return an error to the client.", func() {
					handlers.GetRepository(mClient, "some stuff that looks like an invalid json")
					<-mClient.FinishedChannels[RepositoryFinished]
					Expect(mClient.Name).To(ContainSubstring(expectedName))
					Expect(mClient.Data).To(ContainSubstring(expectedData))
				})
			})
			Context("and repositories are not retrieved", func() {
				expectedName := "logout"
				expectedData := "Cannot get the user repositories."
				It("should return an error to the client.", func() {
					mockWithValidToken := &MockRepositories{Names: []string{""}, AccessToken: "90d64460d14870c08c81352a05dedd3465940a7c"}
					mockJsonWithValidToken := structs.Map(mockWithValidToken)
					mux = http.NewServeMux()
					testServer = httptest.NewServer(mux)
					defer testServer.Close()
					mux.HandleFunc("/user/repos/", func(w http.ResponseWriter, r *http.Request) {
						fmt.Fprint(w, `{"Name":["repo1", "repo2", "repo3"]}`)
					})
					// To simulate the failed test a blank url is given to the
					// github client server. This way, the server can't get any // username and instead returns an error.
					mClient.OauthConf = &oauth2.Config{
						ClientID:     "FAKE_CLIENT_ID",
						ClientSecret: "FAKE_CLIENT_SECRET",
						Endpoint: oauth2.Endpoint{
							TokenURL: "",
						},
					}
					handlers.GetRepository(mClient, mockJsonWithValidToken)
					<-mClient.FinishedChannels[RepositoryFinished]
					Expect(mClient.Name).To(ContainSubstring(expectedName))
					Expect(mClient.Data).To(Equal(expectedData))
				})
			})
			Context("and repositories are successfully retrieved", func() {
				expectedName := "repository set"
				expectedData := &MockRepositories{Names: []string{"repo1", "repo2", "repo3"}, AccessToken: "90d64460d14870c08c81352a05dedd3465940a7c"}
				It("should return the retrieved repositories to the client.", func() {
					mockWithValidToken := &MockRepositories{Names: []string{""}, AccessToken: "90d64460d14870c08c81352a05dedd3465940a7c"}
					mockJsonWithValidToken := structs.Map(mockWithValidToken)
					mux = http.NewServeMux()
					testServer = httptest.NewServer(mux)
					defer testServer.Close()
					mux.HandleFunc("/user/repos/", func(w http.ResponseWriter, r *http.Request) {
						fmt.Fprint(w, `[{"Name":"repo1"}, {"Name":"repo2"}, {"Name":"repo3"}]`)
					})
					// To simulate the failed test a blank url is given to the
					// github client server. This way, the server can't get any // username and instead returns an error.
					mClient.OauthConf = &oauth2.Config{
						ClientID:     "FAKE_CLIENT_ID",
						ClientSecret: "FAKE_CLIENT_SECRET",
						Endpoint: oauth2.Endpoint{
							TokenURL: testServer.URL,
						},
					}
					handlers.GetRepository(mClient, mockJsonWithValidToken)
					<-mClient.FinishedChannels[RepositoryFinished]
					receivedData := mClient.Data.(handlers.Repositories)
					Expect(mClient.Name).To(ContainSubstring(expectedName))
					Expect(receivedData.Names).To(Equal(expectedData.Names))
					Expect(receivedData.AccessToken).To(Equal(expectedData.AccessToken))
				})
			})
		})
	})
})

package handlers

import (
	"github.com/gorilla/websocket"
	models "github.com/joaodias/hugito-app/models"
	"golang.org/x/oauth2"
	githuboauth "golang.org/x/oauth2/github"
	"os"
)

// FindHandler returns the Handler related to the given message sent by the client.
type FindHandler func(string) (Handler, bool)

// Message represents the structure of the messages exchanged between the client and the server.
type Message struct {
	Name string      `json:"name"`
	Data interface{} `json:"data"`
}

// Communicator consists of a set of functions related to the communication and
// authentication process. Communicator is passed to all the handlers so they
// can use these basic functionalities.
type Communicator interface {
	Read()
	Write()
	SetSend(string, interface{})
	GetOauthConfiguration() *oauth2.Config
	GetDBSession() models.DataStorage
	GithubWrapper
}

// SocketClient is a structure to make easier the communication with the client.
type SocketClient struct {
	send        chan Message
	socket      *websocket.Conn
	findHandler FindHandler
	oauthConf   *oauth2.Config
	session     *DBSession
}

// Read reads from the websocket.
func (socketClient *SocketClient) Read() {
	var message Message
	for {
		if err := socketClient.socket.ReadJSON(&message); err != nil {
			break
		}
		if handler, found := socketClient.findHandler(message.Name); found {
			handler(socketClient, message.Data)
		}
	}
	socketClient.socket.Close()
}

// Write writes to the websocket.
func (socketClient *SocketClient) Write() {
	for msg := range socketClient.send {
		if err := socketClient.socket.WriteJSON(msg); err != nil {
			break
		}
	}
	socketClient.socket.Close()
}

// SetSend sends a message through the websocket. The name represents the
// operation and data is the message to send.
func (socketClient *SocketClient) SetSend(name string, data interface{}) {
	socketClient.send <- Message{name, data}
}

// GetOauthConfiguration gets the oauth configuration of a given communicator.
func (socketClient *SocketClient) GetOauthConfiguration() *oauth2.Config {
	return socketClient.oauthConf
}

// GetDBSession gets the database session for the client.
func (socketClient *SocketClient) GetDBSession() models.DataStorage {
	return socketClient.session
}

// NewClient creates a new client communication
func NewClient(socket *websocket.Conn, findHandler FindHandler, session *DBSession) Communicator {
	return &SocketClient{
		send:        make(chan Message),
		socket:      socket,
		findHandler: findHandler,
		oauthConf: &oauth2.Config{
			ClientID:     os.Getenv("CLIENTID"),
			ClientSecret: os.Getenv("SECRET"),
			Endpoint:     githuboauth.Endpoint,
		},
		session: session,
	}
}

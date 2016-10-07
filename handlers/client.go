package handlers

import (
	"github.com/gorilla/websocket"
	"golang.org/x/oauth2"
	githuboauth "golang.org/x/oauth2/github"
)

// Consts used to numerically identify each handler. Used with channels to
// check the related handler.
const (
	UserFinished = iota
	AuthenticationFinished
	RepositoryFinished
	ContentFinished
	ValidationFinished
	FileContentFinished
	PublishContentFinished
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
	NewFinishedChannel(int)
	FinishForKey(int)
	Finished(int)
	GetOauthConfiguration() *oauth2.Config
	GithubWrapper
}

// SocketClient is a structure to make easier the communication with the client.
type SocketClient struct {
	send             chan Message
	socket           *websocket.Conn
	findHandler      FindHandler
	finishedChannels map[int]chan bool
	oauthConf        *oauth2.Config
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

// NewFinishedChannel makes a new channel related to an handler. This channel
// turns true whenever the operation related to that handler is finished.
func (socketClient *SocketClient) NewFinishedChannel(finishedKey int) {
	socketClient.FinishForKey(finishedKey)
	socketClient.finishedChannels[finishedKey] = make(chan bool)
}

// FinishForKey finishes and deletes the channel associated to a given key in
// the case that the channel already exists. This key numerically represents an
// handler.
func (socketClient *SocketClient) FinishForKey(key int) {
	if ch, found := socketClient.finishedChannels[key]; found {
		ch <- true
		delete(socketClient.finishedChannels, key)
	}
}

// Finished sets a specific channel from the set of channels that represent the
// end of an operation to true. This value represents that the operation
// associated to the given channel is finished.
func (socketClient *SocketClient) Finished(key int) {
	socketClient.finishedChannels[key] <- true
}

// GetOauthConfiguration gets the oauth configuration of a given communicator.
func (socketClient *SocketClient) GetOauthConfiguration() *oauth2.Config {
	return socketClient.oauthConf
}

// NewClient creates a new client communication
func NewClient(socket *websocket.Conn, findHandler FindHandler) Communicator {
	return &SocketClient{
		send:             make(chan Message),
		socket:           socket,
		findHandler:      findHandler,
		finishedChannels: make(map[int]chan bool),
		oauthConf: &oauth2.Config{
			ClientID:     ClientID,
			ClientSecret: Secret,
			Endpoint:     githuboauth.Endpoint,
		},
	}
}

package handlers

import (
	"fmt"
	"github.com/gorilla/websocket"
	models "github.com/joaodias/hugito-app/models"
	"net/http"
)

// DBSession wraps the database session.
type DBSession struct {
	*models.Session
}

// Handler is a function that represents the handlers used to handle the messages received by the client.
type Handler func(Communicator, interface{})

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// Router holds the mapping between messages and the respective handlers and
// also the database session.
type Router struct {
	rules   map[string]Handler
	session *DBSession
}

// NewRouter creates a new router with the mapping of the messages to the respective handlers and the session of the database.
func NewRouter(session *DBSession) *Router {
	return &Router{
		rules:   make(map[string]Handler),
		session: session,
	}
}

// Handle is called to assign the defined handler to the message received by the client.
func (e *Router) Handle(msgName string, handler Handler) {
	e.rules[msgName] = handler
}

// FindHandler gives the handler associated to a message.
func (e *Router) FindHandler(msgName string) (Handler, bool) {
	handler, found := e.rules[msgName]
	return handler, found
}

// ServeHTTP upgrades the communication to support websockets and initiallizes the connection. It implements the Handler interface.
func (e *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	socket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
		return
	}
	communicator := NewClient(socket, e.FindHandler, e.session)
	go communicator.Write()
	communicator.Read()
}

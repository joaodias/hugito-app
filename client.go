package main

import (
    "github.com/gorilla/websocket"
)

// FindHandler returns the Handler related to the given message sent by the client.
type FindHandler func(string) (Handler, bool)

// Message represents the structure of the messages exchanged between the client and the server.
type Message struct {
    Name string      `json:"name"`
    Data interface{} `json:"data"`
}

// Client is a structure to make easier the communication with the client.
type Client struct {
    send        chan Message
    socket      *websocket.Conn
    findHandler FindHandler
}

func (client *Client) Read() {
    var message Message
    for {
        if err := client.socket.ReadJSON(&message); err != nil {
            break
        }
        if handler, found := client.findHandler(message.Name); found {
            handler(client, message.Data)
        }
    }
    client.socket.Close()
}

func (client *Client) Write() {
    for msg := range client.send {
        if err := client.socket.WriteJSON(msg); err != nil {
            break
        }
    }
    client.socket.Close()
}

// NewClient creates a new client communication
func NewClient(socket *websocket.Conn, findHandler FindHandler) *Client {
    return &Client{
        send:        make(chan Message),
        socket:      socket,
        findHandler: findHandler,
    }
}

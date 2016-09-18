package handlers

import (
    "fmt"
)

// Content represents the information exchanged between the server and the client.
type Content struct {
    title  string
    author string
    date   string
}

func SubscribeContent(communicator Communicator, data interface{}) {
    fmt.Print("Subscribe content \n")
}

func UnsubscribeContent(communicator Communicator, data interface{}) {
    fmt.Print("Unsubscribe content \n")
}

func AddContent(communicator Communicator, data interface{}) {
    fmt.Print("Add content \n")
}

func RemoveContent(communicator Communicator, data interface{}) {
    fmt.Print("Remove content \n")
}

func UpdateContent(communicator Communicator, data interface{}) {
    fmt.Print("Update content \n")
}

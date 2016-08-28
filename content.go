package main

import (
    "fmt"
)

// Content represents the information exchanged between the server and the client.
type Content struct {
    title  string
    author string
    date   string
}

func subscribeContent(client *Client, data interface{}) {
    fmt.Print("Subscribe content \n")
}

func unsubscribeContent(client *Client, data interface{}) {
    fmt.Print("Unsubscribe content \n")
}

func addContent(client *Client, data interface{}) {
    fmt.Print("Add content \n")
}

func removeContent(client *Client, data interface{}) {
    fmt.Print("Remove content \n")
}

func updateContent(client *Client, data interface{}) {
    fmt.Print("Update content \n")
}

package main

import (
    "fmt"
)

// User represents the user information exchanged between the server and the client.
type User struct {
    Name string
}

func subscribeUser(client *Client, data interface{}) {
    fmt.Print("User subscribe \n")
}

func unsubscribeUser(client *Client, data interface{}) {
    fmt.Print("User unsubscribe \n")
}

func setUser(client *Client, data interface{}) {
    fmt.Print("User set /n")
}

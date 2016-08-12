package main

import (
	"fmt"
)

type Repository struct {
	Name string
}

func subscribeRepository(client *Client, data interface{}) {
	fmt.Print("Repository subscribe \n")
}

func unsubscribeRepository(client *Client, data interface{}) {
	fmt.Print("Repository unsubscribe \n")
}

func addRepository(client *Client, data interface{}) {
	fmt.Print("Add repository \n")
}

func removeRepository(client *Client, data interface{}) {
	fmt.Print("Remove repository \n")
}

func validateRepository(client *Client, data interface{}) {
	fmt.Print("Validate repository \n")

	client.send <- Message{"repository validate", true}
}

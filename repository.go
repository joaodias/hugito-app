package main

import (
	"fmt"
)

type Repository struct {
	Name string
}

func addRepository(client *Client, data interface{}) {
	fmt.Print("Add repository \n")
}

func removeRepository(client *Client, data interface{}) {
	fmt.Print("Remove repository \n")
}

func validateRepository(client *Client, data interface{}) {
	fmt.Print("Validate repository \n")
}

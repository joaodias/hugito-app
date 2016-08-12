package main

import (
	"fmt"
)

type Content struct {
	title  string
	author string
	date   string
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

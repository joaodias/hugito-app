package main

import (
    "fmt"
)

type Configuration struct {
    FieldNames: string[]
    FieldValues: string[]
}

func subscribeConfiguration(client *Client, data interface{}) {
    fmt.Print("Subscribe configuration \n")
}

func unsubscribeConfiguration(client *Client, data interface{}) {
    fmt.Print("Unsubscribe configuration \n")
}

func addConfiguration(client *Client, data interface{}) {
    fmt.Print("Add configuration \n")
}

func removeConfiguration(client *Client, data interface{}) {
    fmt.Print("Remove configuration \n")
}

func updateConfiguration(client *Client, data interface{}) {
    fmt.Print("Update configuration \n")
}

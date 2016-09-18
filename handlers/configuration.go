package handlers

import (
    "fmt"
)

// Configuration holds the values exchanged between the client and the server related to the website configuration file.
type Configuration struct {
    FieldNames  []string
    FieldValues []string
}

func SubscribeConfiguration(communicator Communicator, data interface{}) {
    fmt.Print("Subscribe configuration \n")
}

func UnsubscribeConfiguration(communicator Communicator, data interface{}) {
    fmt.Print("Unsubscribe configuration \n")
}

func AddConfiguration(communicator Communicator, data interface{}) {
    fmt.Print("Add configuration \n")
}

func RemoveConfiguration(communicator Communicator, data interface{}) {
    fmt.Print("Remove configuration \n")
}

func UpdateConfiguration(communicator Communicator, data interface{}) {
    fmt.Print("Update configuration \n")
}

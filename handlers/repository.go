package handlers

import (
    "fmt"
)

// Repository represents the repository parameters exchanged between the server and the client.
type Repository struct {
    Name string
}

func SubscribeRepository(communicator Communicator, data interface{}) {
    fmt.Print("Repository subscribe \n")
}

func UnsubscribeRepository(communicator Communicator, data interface{}) {
    fmt.Print("Repository unsubscribe \n")
}

func AddRepository(communicator Communicator, data interface{}) {
    fmt.Print("Add repository \n")
}

func RemoveRepository(communicator Communicator, data interface{}) {
    fmt.Print("Remove repository \n")
}

func ValidateRepository(communicator Communicator, data interface{}) {
    fmt.Print("Validate repository \n")
}

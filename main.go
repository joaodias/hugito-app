package main

import (
	"fmt"
	handlers "github.com/joaodias/hugito-app/handlers"
	"net/http"
)

func main() {
	router := handlers.NewRouter()
	router.Handle("repository get", handlers.GetRepository)
	router.Handle("repository validate", handlers.ValidateRepository)
	router.Handle("content list", handlers.GetContentList)
	router.Handle("user get", handlers.GetUser)
	router.Handle("authenticate", handlers.Authenticate)
	http.Handle("/", router)
	fmt.Print("Go app initialized on port 4000 \n")
	http.ListenAndServe(":4000", nil)
}

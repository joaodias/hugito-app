package main

import (
	"fmt"
	"github.com/joaodias/hugito-app/handlers"
	"net/http"
	"os"
)

func main() {
	router := handlers.NewRouter()
	router.Handle("repositories get", handlers.GetRepository)
	router.Handle("repository validate", handlers.ValidateRepository)
	router.Handle("content list", handlers.GetContentList)
	router.Handle("content get", handlers.GetFileContent)
	router.Handle("content create", handlers.CreateContent)
	router.Handle("content update", handlers.UpdateContent)
	router.Handle("content remove", handlers.RemoveContent)
	router.Handle("user get", handlers.GetUser)
	router.Handle("authenticate", handlers.Authenticate)
	http.Handle("/", router)
	port := os.Getenv("PORT")
	fmt.Print("Go app initialized in port " + port + ".\n")
	http.ListenAndServe(":"+port, nil)
}

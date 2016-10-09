package main

import (
	log "github.com/Sirupsen/logrus"
	handlers "github.com/joaodias/hugito-app/handlers"
	"net/http"
)

func main() {
	router := handlers.NewRouter()
	router.Handle("repositories get", handlers.GetRepository)
	router.Handle("repository validate", handlers.ValidateRepository)
	router.Handle("content list", handlers.GetContentList)
	router.Handle("content get", handlers.GetFileContent)
	router.Handle("content update", handlers.PublishContent)
	router.Handle("user get", handlers.GetUser)
	router.Handle("authenticate", handlers.Authenticate)
	http.Handle("/", router)
	log.Info("Go app initialized in port 4000.")
	http.ListenAndServe(":4000", nil)
}

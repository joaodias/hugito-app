package main

import (
	"fmt"
	"github.com/joaodias/hugito-app/handlers"
	models "github.com/joaodias/hugito-app/models"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func main() {
	// This needs to be changed depending on the environment you are working.
	err := godotenv.Load("production.env")
	if err != nil {
		log.Panic("Error loading .env file")
	}
	session, err := models.InitSession()
	if err != nil {
		log.Panic("Error initializing session: " + err.Error())
	}
	router := handlers.NewRouter(&handlers.DBSession{session})
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

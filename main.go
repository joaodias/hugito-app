package main

import (
	"fmt"
	"net/http"

	handlers "github.com/joaodias/hugito-app/handlers"
)

func main() {
	router := handlers.NewRouter()

	router.Handle("repository get", handlers.GetRepository)
	router.Handle("repository unsubscribe", handlers.UnsubscribeRepository)
	router.Handle("repository add", handlers.AddRepository)
	router.Handle("repository remove", handlers.RemoveRepository)
	router.Handle("repository validate", handlers.ValidateRepository)

	router.Handle("content subscribe", handlers.SubscribeContent)
	router.Handle("content unsubscribe", handlers.UnsubscribeContent)
	router.Handle("content add", handlers.AddContent)
	router.Handle("content remove", handlers.RemoveContent)
	router.Handle("content update", handlers.UpdateContent)

	router.Handle("configuration subscribe", handlers.SubscribeContent)
	router.Handle("configuration unsubscribe", handlers.UnsubscribeContent)
	router.Handle("configuration add", handlers.AddContent)
	router.Handle("configuration remove", handlers.RemoveContent)
	router.Handle("configuration update", handlers.UpdateContent)

	router.Handle("user get", handlers.GetUser)

	router.Handle("authenticate", handlers.Authenticate)

	http.Handle("/", router)

	fmt.Print("Go app initialized on port 4000 \n")

	http.ListenAndServe(":4000", nil)
}

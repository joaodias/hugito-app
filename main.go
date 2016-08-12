package main

import (
	"fmt"
	"net/http"
)

func main() {
	router := NewRouter()

	router.Handle("repository add", addRepository)
	router.Handle("repository remove", removeRepository)
	router.Handle("repository validate", validateRepository)

	router.Handle("content add", addContent)
	router.Handle("content remove", removeContent)
	router.Handle("content update", updateContent)

	router.Handle("user subscribe", subscribeUser)
	router.Handle("user unsubscribe", unsubscribeUser)
	router.Handle("user set", setUser)

	http.Handle("/", router)

	fmt.Print("Go app initialized on port 4000 \n")

	http.ListenAndServe(":4000", nil)
}

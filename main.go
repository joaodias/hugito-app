package main

import (
	"github.com/gorilla/mux"
	"github.com/joaodias/hugito-backend/infrastructure"
	"github.com/joaodias/hugito-backend/interfaces"
	"github.com/joaodias/hugito-backend/usecases"
	"golang.org/x/oauth2"
	githuboauth "golang.org/x/oauth2/github"
	"net/http"
	"os"
)

func main() {
	dbHandler, err := infrastructure.NewRethinkDBSession()
	if err != nil {
		panic("Can't connect to the database")
	}
	sourceControl := &infrastructure.Github{}

	buildEngine := &infrastructure.Hugo{}

	userInteractor := usecases.UserInteractor{}
	userInteractor.UserRepository = &interfaces.ExternalUserRepository{
		SourceControl:   sourceControl,
		DatabaseHandler: dbHandler,
		BuildEngine:     buildEngine,
		Logger:          &infrastructure.Logger{},
	}

	contentInteractor := usecases.ContentInteractor{}
	contentInteractor.ContentRepository = &interfaces.ExternalContentRepository{
		SourceControl:   sourceControl,
		DatabaseHandler: dbHandler,
		BuildEngine:     buildEngine,
		Logger:          &infrastructure.Logger{},
	}

	repositoryInteractor := usecases.RepositoryInteractor{}
	repositoryInteractor.RepositoryRepository = &interfaces.ExternalRepositoryRepository{
		SourceControl:   sourceControl,
		DatabaseHandler: dbHandler,
		Logger:          &infrastructure.Logger{},
	}

	webHandler := &interfaces.WebHandler{
		OauthConfiguration: &oauth2.Config{
			ClientID:     os.Getenv("CLIENTID"),
			ClientSecret: os.Getenv("SECRET"),
			Endpoint:     githuboauth.Endpoint,
		},
		UserInteractor:       &userInteractor,
		ContentInteractor:    &contentInteractor,
		RepositoryInteractor: &repositoryInteractor,
		Logger:               &infrastructure.Logger{},
	}

	router := mux.NewRouter()
	router.HandleFunc("/user", webHandler.ReadUser).Methods("POST")
	router.HandleFunc("/auth", webHandler.Authenticate).Methods("POST")
	router.HandleFunc("/auth", webHandler.AuthWebOptions).Methods("OPTIONS")
	router.HandleFunc("/content", webHandler.CreateContent).Methods("POST")
	router.HandleFunc("/content", webHandler.ContentWebOptions).Methods("OPTIONS")
	router.HandleFunc("/content", webHandler.ListContent).Methods("GET")
	router.HandleFunc("/content/{name}", webHandler.ReadContent).Methods("GET")
	router.HandleFunc("/content", webHandler.SaveContent).Methods("PUT")
	router.HandleFunc("/content", webHandler.RemoveContent).Methods("DELETE")
	router.HandleFunc("/content/publish", webHandler.ContentWebOptions).Methods("OPTIONS")
	router.HandleFunc("/content/publish", webHandler.PublishContent).Methods("POST")
	router.HandleFunc("/repository", webHandler.ValidateRepository).Methods("GET")
	port := os.Getenv("PORT")
	logger := &infrastructure.Logger{}
	logger.Log("Server running on port " + port)
	http.ListenAndServe(port, router)
}

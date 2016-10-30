package models

import (
	r "github.com/dancannon/gorethink"
	"os"
)

// DataStorage holds the the methods related to data storage situations.
type DataStorage interface {
	AddUser(interface{}) error
}

// Session embeds the database session logic.
type Session struct {
	*r.Session
}

// InitSession initiallizes the database session.
func InitSession() (*Session, error) {
	session, err := r.Connect(r.ConnectOpts{
		Address:  os.Getenv("DBHOST"),
		Database: os.Getenv("DBNAME"),
	})
	if err != nil {
		return nil, err
	}
	return &Session{session}, nil
}

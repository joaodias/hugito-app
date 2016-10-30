package models

import (
	r "github.com/dancannon/gorethink"
	"os"
)

// Host server and database name.
var (
	DBHost = os.Getenv("DBHOST")
	DBName = os.Getenv("DBNAME")
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
		Address:  DBHost,
		Database: DBName,
	})
	if err != nil {
		return nil, err
	}
	return &Session{session}, nil
}

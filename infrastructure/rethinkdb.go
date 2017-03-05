package infrastructure

import (
	r "github.com/dancannon/gorethink"
	"github.com/joaodias/hugito-backend/domain"
	"os"
)

// RethinkDBSession holds the rethinkDB session
type RethinkDBSession struct {
	*r.Session
}

// Add adds data to a given table. Also a conflict solver is passed.
func (rs *RethinkDBSession) Add(data interface{}, table string, onConflict string) error {
	err := r.Table(table).
		Insert(data, r.InsertOpts{
			Conflict: onConflict,
		}).
		Exec(rs.Session)
	return err
}

// Remove removes an entry with id from a given table
func (rs *RethinkDBSession) Remove(id string, table string) error {
	err := r.Table(table).Get(id).Delete().
		Exec(rs.Session)
	return err
}

// Update updates an entry with id from a given table
func (rs *RethinkDBSession) Update(data interface{}, id string, table string) error {
	err := r.Table(table).Get(id).Update(data).
		Exec(rs.Session)
	return err
}

// List list a database table given given a key field
func (rs *RethinkDBSession) List(key string, table string) ([]domain.Content, error) {
	contents := make([]domain.Content, 0)
	res, err := r.Table(table).GetAllByIndex(key).Run(rs.Session)
	res.All(&contents)
	return contents, err
}

// NewRethinkDBSession creates a new database session
func NewRethinkDBSession() (*RethinkDBSession, error) {
	session, err := r.Connect(r.ConnectOpts{
		Address:  os.Getenv("DBHOST"),
		Database: os.Getenv("DBNAME"),
	})
	if err != nil {
		return nil, err
	}
	return &RethinkDBSession{session}, nil
}

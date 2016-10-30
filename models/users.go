package models

import (
	r "github.com/dancannon/gorethink"
)

// AddUser adds an user to the database. If the user exists, updates it.
func (session *Session) AddUser(user interface{}) error {
	err := r.Table("users").
		Insert(user, r.InsertOpts{
			Conflict: "update",
		}).
		Exec(session)
	return err
}

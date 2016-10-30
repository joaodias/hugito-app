package mocks

import (
	"errors"
	models "github.com/joaodias/hugito-app/models"
)

// Mocking for the DataStorage.

// The flag isError is used to control the return of the methods.
type DataStorage struct {
	IsError bool
}

var _ models.DataStorage = (*DataStorage)(nil)

func (dataStorage *DataStorage) AddUser(mockUser interface{}) error {
	if dataStorage.IsError {
		return errors.New("Some error.")
	}
	return nil
}

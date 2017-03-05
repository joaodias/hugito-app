package mocks

import (
	"errors"
	"github.com/joaodias/hugito-backend/domain"
)

type DatabaseHandler struct {
	IsError      bool
	AddCalled    bool
	RemoveCalled bool
	UpdateCalled bool
	ListCalled   bool
	Table        string
	Conflict     string
}

func (dh *DatabaseHandler) Add(data interface{}, table string, onConflict string) error {
	dh.AddCalled = true
	dh.Table = table
	dh.Conflict = onConflict
	if dh.IsError {
		return errors.New("Some error")
	}
	return nil
}

func (dh *DatabaseHandler) Remove(id string, table string) error {
	dh.RemoveCalled = true
	dh.Table = table
	if dh.IsError {
		return errors.New("Some error")
	}
	return nil
}

func (dh *DatabaseHandler) Update(data interface{}, id string, table string) error {
	dh.UpdateCalled = true
	dh.Table = table
	if dh.IsError {
		return errors.New("Some error")
	}
	return nil
}

func (dh *DatabaseHandler) List(id string, table string) ([]domain.Content, error) {
	dh.ListCalled = true
	dh.Table = table
	if dh.IsError {
		return nil, errors.New("Some error")
	}
	return make([]domain.Content, 1), nil
}

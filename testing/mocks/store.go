package mocks

import (
	"testing"
)

type Store struct {
	GetFunc       func(key string) (string, error)
	SetFunc       func(key string, value string) error
	GetRecordFunc func(key string, value interface{}) error
	SetRecordFunc func(key string, value interface{}) error
}

func BaselineStore(t *testing.T) *Store {
	t.Helper()

	store := Store{
		GetFunc: func(string) (string, error) {
			return GenericString, nil
		},
		SetFunc: func(string, string) error {
			return nil
		},
		GetRecordFunc: func(string, interface{}) error {
			return nil
		},
		SetRecordFunc: func(string, interface{}) error {
			return nil
		},
	}

	return &store
}

func (s *Store) Get(key string) (string, error) {
	return s.GetFunc(key)
}

func (s *Store) Set(key string, value string) error {
	return s.SetFunc(key, value)
}

func (s *Store) GetRecord(key string, value interface{}) error {
	return s.GetRecordFunc(key, value)
}

func (s *Store) SetRecord(key string, value interface{}) error {
	return s.SetRecordFunc(key, value)
}

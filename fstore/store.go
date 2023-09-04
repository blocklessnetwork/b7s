package fstore

type Store interface {
	GetRecord(string, interface{}) error
	SetRecord(string, interface{}) error
	Keys() ([]string, error)
}

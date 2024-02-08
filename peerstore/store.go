package peerstore

type Store interface {
	SetRecord(key string, value interface{}) error
	GetRecord(key string, out interface{}) error
	Keys() []string
	Delete(key string) error
}

package function

type Store interface {
	GetRecord(string, interface{}) error
	SetRecord(string, interface{}) error
}

package fstore

import (
	"github.com/blocklessnetwork/b7s/models/blockless"
)

type Store interface {
	RetrieveFunction(cid string) (blockless.FunctionRecord, error)
	SaveFunction(cid string, function blockless.FunctionRecord) error
	RetrieveFunctions() ([]blockless.FunctionRecord, error)
}

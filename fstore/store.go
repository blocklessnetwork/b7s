package fstore

import (
	"github.com/blocklessnetwork/b7s/models/blockless"
)

type Store interface {
	RetrieveFunction(cid string) (blockless.FunctionRecord, error)
	SaveFunction(function blockless.FunctionRecord) error
	RetrieveFunctions() ([]blockless.FunctionRecord, error)
}

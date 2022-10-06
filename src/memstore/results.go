package memstore

import (
	"errors"
	"sync"

	"github.com/blocklessnetworking/b7s/src/models"
	log "github.com/sirupsen/logrus"
)

type ReqRespStore interface {
	Get(string) *models.MsgExecuteResponse
	Set(string, *models.MsgExecuteResponse) error
}

type respDB map[string]*models.MsgExecuteResponse

type localStorageReqRespStore struct {
	db   respDB
	lock sync.RWMutex
}

func (rs *localStorageReqRespStore) Get(requestId string) *models.MsgExecuteResponse {
	rs.lock.RLock()
	resp, ok := rs.db[requestId]
	rs.lock.RUnlock()
	if !ok {
		return nil
	}

	return resp
}

func (rs *localStorageReqRespStore) Set(requestId string, resp *models.MsgExecuteResponse) error {
	rs.lock.Lock()
	defer rs.lock.Unlock()
	existing, ok := rs.db[requestId]
	if !ok {
		rs.db[requestId] = resp
		return nil
	} else {
		log.WithFields(log.Fields{
			"requestId": requestId,
			"resp":      resp,
			"existing":  existing,
		}).Error("Response already exists")
		return errors.New("Response for request " + requestId + " already exists")
	}
}

func NewReqRespStore() ReqRespStore {
	return &localStorageReqRespStore{
		db: make(respDB),
	}
}

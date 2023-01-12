package memstore

import (
	"testing"

	"github.com/blocklessnetworking/b7s/src/enums"
	"github.com/blocklessnetworking/b7s/src/models"
	"github.com/stretchr/testify/assert"
)

func TestReqRespStore_Set(t *testing.T) {
	store := NewReqRespStore()

	// Test inserting a new request-response pair
	requestID := "123"
	response := &models.MsgExecuteResponse{
		Code:   enums.ResponseCodeOk,
		Result: "hello world",
	}
	err := store.Set(requestID, response)
	assert.Nil(t, err)

	// Test inserting a duplicate request-response pair
	err = store.Set(requestID, response)
	assert.NotNil(t, err)
	assert.Equal(t, "Response for request 123 already exists", err.Error())

	// Test getting the inserted request-response pair
	storedResponse := store.Get(requestID)
	assert.Equal(t, response, storedResponse)
}

package db

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDb(t *testing.T) {
	// setup
	databaseID := "/tmp/test_db"
	os.RemoveAll(databaseID)

	// test GetDb
	db := GetDb(databaseID)
	assert.NotNil(t, db)

	ctx := context.WithValue(context.Background(), "appDb", db)

	// test Set and Get
	err := Set(ctx, "test_key", "test_value")
	assert.Nil(t, err)

	value, err := Get(ctx, "test_key")
	assert.Nil(t, err)
	assert.Equal(t, "test_value", string(value))

	// test GetString
	stringValue, err := GetString(ctx, "test_key")
	assert.Nil(t, err)
	assert.Equal(t, "test_value", stringValue)

	// test Close
	Close(ctx)
	value, err = Get(ctx, "test_key")
	assert.Nil(t, value)
	assert.NotNil(t, err)
}

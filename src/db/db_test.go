package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDb(t *testing.T) {
	// setup
	appDb := GetDb("/tmp/test_db")
	defer Close(appDb)
	ctx := context.WithValue(context.Background(), "appDb", appDb)

	// test set
	err := Set(ctx, "test_key", "test_value")
	assert.Nil(t, err)

	// test get
	val, err := Get(ctx, "test_key")
	assert.Nil(t, err)
	assert.Equal(t, "test_value", string(val))

	// test get string
	valStr, err := GetString(ctx, "test_key")
	assert.Nil(t, err)
	assert.Equal(t, "test_value", valStr)
}

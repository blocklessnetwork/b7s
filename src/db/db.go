package db

import (
	"context"
	"fmt"
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/cockroachdb/pebble"
)

var (
	mtx      sync.RWMutex
	db       *pebble.DB
	isClosed bool
)

func GetDb(databaseID string) *pebble.DB {
	mtx.Lock()
	defer mtx.Unlock()
	if db == nil {
		log.Info("opening database: ", databaseID)
		d, err := pebble.Open(databaseID, &pebble.Options{})
		if err != nil {
			log.Warn(err)
		}
		db = d
	}
	isClosed = false
	return db
}

func Set(ctx context.Context, key string, value string) error {
	mtx.Lock()
	defer mtx.Unlock()

	if isClosed {
		return fmt.Errorf("database is closed")
	}

	d := ctx.Value("appDb").(*pebble.DB)
	if err := d.Set([]byte(key), []byte(value), pebble.Sync); err != nil {
		log.Warn(err)
		return err
	}
	return nil
}

func Get(ctx context.Context, key string) ([]byte, error) {
	mtx.RLock()
	defer mtx.RUnlock()

	if isClosed {
		return nil, fmt.Errorf("database is closed")
	}

	d := ctx.Value("appDb").(*pebble.DB)
	value, closer, err := d.Get([]byte(key))
	if err != nil {
		return nil, err
	}
	defer closer.Close()
	return value, nil
}

func GetString(ctx context.Context, key string) (string, error) {
	mtx.RLock()
	defer mtx.RUnlock()

	if isClosed {
		return "", fmt.Errorf("database is closed")
	}

	d := ctx.Value("appDb").(*pebble.DB)
	value, closer, err := d.Get([]byte(key))
	if err != nil {
		return "", err
	}
	stringVal := string(value)
	defer closer.Close()
	return stringVal, nil
}

func Close(ctx context.Context) error {
	mtx.Lock()
	defer mtx.Unlock()

	if isClosed {
		return fmt.Errorf("database is closed")
	}

	d := ctx.Value("appDb").(*pebble.DB)
	if db != nil {
		if err := d.Close(); err != nil {
			log.Warn(err)
		}
		db = nil
	}
	isClosed = true
	return nil
}

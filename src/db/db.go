package db

import (
	"context"

	log "github.com/sirupsen/logrus"

	"github.com/cockroachdb/pebble"
)

func GetDb(DatabaseId string) *pebble.DB {
	dbPath := DatabaseId
	log.Info("opening database: ", dbPath)
	db, err := pebble.Open(dbPath, &pebble.Options{})
	if err != nil {
		log.Warn(err)
	}
	return db
}

func Set(ctx context.Context, key string, value string) error {
	db := ctx.Value("appDb").(*pebble.DB)
	if err := db.Set([]byte(key), []byte(value), pebble.Sync); err != nil {
		log.Warn(err)
		return err
	}
	return nil
}

func Get(ctx context.Context, key string) ([]byte, error) {
	db := ctx.Value("appDb").(*pebble.DB)
	value, closer, err := db.Get([]byte(key))
	if err != nil {
		return nil, err
	}
	defer closer.Close()
	return value, nil
}

func GetString(ctx context.Context, key string) (string, error) {
	db := ctx.Value("appDb").(*pebble.DB)
	value, closer, err := db.Get([]byte(key))
	if err != nil {
		return "", err
	}
	stringVal := string(value)
	defer closer.Close()
	return stringVal, nil
}

func Close(db *pebble.DB) {
	if err := db.Close(); err != nil {
		log.Warn(err)
	}
}

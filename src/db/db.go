package db

import (
	log "github.com/sirupsen/logrus"

	"github.com/cockroachdb/pebble"
)

func Get(DatabaseId string) *pebble.DB {

	dbPath := DatabaseId
	log.Info("Opening database: ", dbPath)
	db, err := pebble.Open(dbPath, &pebble.Options{})
	if err != nil {
		log.Warn(err)
	}
	return db
}

func Set(db *pebble.DB, key string, value string) error {
	if err := db.Set([]byte(key), []byte(value), pebble.Sync); err != nil {
		log.Warn(err)
		return err
	}
	return nil
}

func Value(db *pebble.DB, key string) string {
	value, closer, err := db.Get([]byte(key))
	if err != nil {
		log.Warn(err)
		return ""
	}
	defer closer.Close()
	return string(value)
}

func Close(db *pebble.DB) {
	if err := db.Close(); err != nil {
		log.Warn(err)
	}
}

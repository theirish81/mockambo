package db

import (
	"github.com/dgraph-io/badger/v4"
)

// DBs is a map of DB connections to badger. We use a map of connections because multiple storages can be used within
// the same OpenAPI file
type DBs map[string]*badger.DB

// databases is the instance of DBs
var databases *DBs

// getDB returns the DB connection for the given storage path. If the connection has been established yet, then
// it gets reused. It gets created otherwise.
func getDB(path string) (*badger.DB, error) {
	if databases == nil {
		databases = &DBs{}
	}
	if db, ok := (*databases)[path]; !ok {
		db, err := badger.Open(badger.DefaultOptions(path).WithLoggingLevel(3))
		if err != nil {
			return nil, err
		}
		(*databases)[path] = db
		return db, nil
	} else {
		return db, nil
	}
}

// Upsert inserts or updates an entry to one of the DBs, identified by the path
func Upsert(key string, value []byte, path string) error {
	db, err := getDB(path)
	if err != nil {
		return err
	}
	return db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(key), value)
	})
}

// Get retrieves an entry from the DB identified by the path. It errors in case the entry is not found
func Get(key string, path string) ([]byte, error) {
	db, err := getDB(path)
	if err != nil {
		return nil, err
	}
	res := make([]byte, 0)
	err = db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			res = val
			return nil
		})
	})
	return res, err
}

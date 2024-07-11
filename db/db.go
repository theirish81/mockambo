package db

import (
	"github.com/dgraph-io/badger/v4"
)

type DBs map[string]*badger.DB

var databases *DBs

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

func Upsert(key string, value []byte, path string) error {
	db, err := getDB(path)
	if err != nil {
		return err
	}
	return db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(key), value)
	})
}

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

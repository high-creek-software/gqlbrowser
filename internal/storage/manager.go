package storage

import (
	"errors"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"os"
	"path/filepath"
)

const (
	appName = "gqlbrowser"
	dbName  = "gqlbrowser.db"
)

type Manager interface {
	Store(path, payload string) (Endpoint, error)
	List() ([]Endpoint, error)
}

type ManagerImpl struct {
	applicationDirectory string

	dbPath string
	db     *gorm.DB

	*endpointRepo
}

func NewManager() Manager {
	m := &ManagerImpl{applicationDirectory: getApplicationDirectory()}

	err := os.Mkdir(m.applicationDirectory, os.ModePerm)
	if err != nil && !errors.Is(err, os.ErrExist) {
		log.Fatal("error creating app directory", err)
	}

	m.dbPath = filepath.Join(m.applicationDirectory, dbName)
	m.db, err = gorm.Open(sqlite.Open(m.dbPath), &gorm.Config{})
	if err != nil {
		log.Fatal("error opening database", err)
	}

	m.endpointRepo = newEndpointRepo(m.db)

	return m
}

package storage

import (
	"encoding/json"
	"gitlab.com/high-creek-software/fieldglass"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"path/filepath"
)

const (
	appName = "gqlbrowser"
	dbName  = "gqlbrowser.db"
)

type Manager interface {
	Create(path string) (Endpoint, error)
	List() ([]Endpoint, error)
	Update(endpoint Endpoint) error
	Delete(endpoint Endpoint) error
}

type ManagerImpl struct {
	applicationDirectory string

	dbPath string
	db     *gorm.DB

	client fieldglass.FieldGlass

	*endpointRepo
}

func NewManager(appDirectory string, client fieldglass.FieldGlass) Manager {
	m := &ManagerImpl{applicationDirectory: appDirectory, client: client}

	//err := os.Mkdir(m.applicationDirectory, os.ModePerm)
	//if err != nil && !errors.Is(err, os.ErrExist) {
	//	log.Fatal("error creating app directory", err)
	//}

	var err error
	m.dbPath = filepath.Join(m.applicationDirectory, dbName)
	m.db, err = gorm.Open(sqlite.Open(m.dbPath), &gorm.Config{})
	if err != nil {
		log.Fatal("error opening database", err)
	}

	m.endpointRepo = newEndpointRepo(m.db)

	return m
}

func (m *ManagerImpl) Create(path string) (Endpoint, error) {
	schema, err := m.client.Load(path)
	if err != nil {
		return Endpoint{}, err
	}

	payload, err := json.Marshal(schema)
	if err != nil {
		return Endpoint{}, err
	}

	endpoint, err := m.Store(path, string(payload))
	if err != nil {
		return Endpoint{}, err
	}

	return endpoint, nil
}

func (m *ManagerImpl) Update(endpoint Endpoint) error {
	schema, err := m.client.Load(endpoint.Path)
	if err != nil {
		return err
	}

	payload, err := json.Marshal(schema)
	if err != nil {
		return err
	}

	return m.endpointRepo.update(endpoint.ID, string(payload))
}

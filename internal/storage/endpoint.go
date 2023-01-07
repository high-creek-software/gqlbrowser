package storage

import (
	"github.com/rs/xid"
	"gorm.io/gorm"
	"log"
	"time"
)

type Endpoint struct {
	ID        string `gorm:"primaryKey"`
	Path      string
	Payload   string
	CreatedAt time.Time
	UpdatedAt *time.Time
}

type endpointRepo struct {
	db *gorm.DB
}

func newEndpointRepo(db *gorm.DB) *endpointRepo {
	err := db.AutoMigrate(&Endpoint{})
	if err != nil {
		log.Fatal("error migrating endpoint model", err)
	}

	return &endpointRepo{db: db}
}

func (e *endpointRepo) Store(path, payload string) (Endpoint, error) {
	ep := Endpoint{ID: xid.New().String(), Path: path, Payload: payload, CreatedAt: time.Now()}

	err := e.db.Create(&ep).Error

	return ep, err
}

func (e *endpointRepo) List() ([]Endpoint, error) {
	var eps []Endpoint
	err := e.db.Find(&eps).Error

	return eps, err
}

func (e *endpointRepo) update(id, schema string) error {
	now := time.Now()
	return e.db.Model(&Endpoint{ID: id}).Updates(Endpoint{Payload: schema, UpdatedAt: &now}).Error
}

func (e *endpointRepo) Delete(endpoint Endpoint) error {
	return e.db.Delete(&endpoint).Error
}

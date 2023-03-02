package repository

import (
	"github.com/google/uuid"
	gormGIS "github.com/nferruzzi/gormgis"
	"ocall/app/entity"
	"time"
)

type JSONB []interface{}

type Application struct {
	Model
	Name         string
	Status       entity.ApplicationStatus `gorm:"type:application_status"`
	PerformerRef uuid.UUID
	EventRef     uuid.UUID
}

type Event struct {
	Model
	Name            string
	Description     string
	Tags            []entity.Tag
	ProducerRef     uuid.UUID
	VenueRef        uuid.UUID
	Applications    []Application `gorm:"foreignKey:EventRef"`
	ApplicationForm JSONB         `gorm:"type:jsonb"`
	ConfirmedActs   []Act         `gorm:"foreignKey:ActRef"`
	Location        gormGIS.GeoPoint
	Status          entity.EventApplicationStatus `gorm:"type:event_application_status"`
	Time            time.Time
	PayStructure    string
}

type Act struct {
	Model
	Name         string
	Tags         []entity.Tag
	Description  string
	Media        JSONB `gorm:"type:jsonb"`
	Form         JSONB `gorm:"type:jsonb"`
	PerformerRef uuid.UUID
	EventRef     uuid.UUID
}

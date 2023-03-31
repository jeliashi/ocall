package models

import (
	"github.com/google/uuid"
	gormGIS "github.com/nferruzzi/gormgis"
	"gorm.io/gorm"
	"time"
)

type Tag struct {
	gorm.Model
	Name string `gorm:"unique"`
}
type GoogleResponseID string
type GoogleFormID string

type ApplicationStatus string

const (
	STATUS_ACCEPTED ApplicationStatus = "accepted"
	STATUS_REJECTED ApplicationStatus = "rejected"
	STATUS_PENDING  ApplicationStatus = "pending"
	STATUS_OFFERED  ApplicationStatus = "offered"
	STATUS_UNKNOWN  ApplicationStatus = "unknown"
)

type EventApplicationStatus string

const (
	EVENT_DRAFT     EventApplicationStatus = "draft"
	EVENT_OPEN      EventApplicationStatus = "open"
	EVENT_CLOSED    EventApplicationStatus = "closed"
	EVENT_CANCELLED EventApplicationStatus = "cancelled"
)

type Application struct {
	Model
	Name             string
	Status           ApplicationStatus `gorm:"type:application_status"`
	Performer        Profile
	EventRef         uuid.UUID
	GoogleResponseID GoogleResponseID
	Saved            bool
}

type Event struct {
	Model
	Name         string
	Description  string
	Tags         []Tag
	Producer     Profile
	Venue        *Profile
	Applications []Application `gorm:"foreignKey:EventRef"`
	GoogleForm   GoogleFormID
	Location     gormGIS.GeoPoint
	Status       EventApplicationStatus `json:"application_status" gorm:"type:event_application_status"`
	Time         time.Time
	ApplyByTime  *time.Time
	PayStructure string
}

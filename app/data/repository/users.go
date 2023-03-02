package repository

import (
	"github.com/google/uuid"
	"github.com/nferruzzi/gormgis"
)

type Performer struct {
	Model
	Name            string        `gorm:"notnull"`
	UserID          uuid.UUID     `gorm:"notnull"`
	Applications    []Application `gorm:"foreignKey:PerformerRef"`
	OfferedEvents   []Event       `gorm:"foreignKey:PerformerRef"`
	ConfirmedEvents []Event       `gorm:"foreignKey:PerformerRef"`
	Acts            []Act         `gorm:"foreignKey:PerformerRef"`
}

type Producer struct {
	Model
	Name      string
	UserID    uuid.UUID
	Events    []Event `gorm:"foreignKey:ProducerRef"`
	SavedActs []Act   `gorm:"foreignKey:ProducerRef"`
}

type Venue struct {
	Model
	Name     string    `json:"name,nonempty"`
	UserID   uuid.UUID `json:"user_id,omitempty"`
	Location gormGIS.GeoPoint
	Events   []Event `gorm:"foreignKey:VenueRef" json:"events,omitempty"`
}

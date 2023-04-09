package models

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/nferruzzi/gormGIS"
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
	StatusAccepted ApplicationStatus = "accepted"
	StatusRejected ApplicationStatus = "rejected"
	StatusPending  ApplicationStatus = "pending"
	StatusOffered  ApplicationStatus = "offered"
	StatusUnknown  ApplicationStatus = "unknown"
)

func (ApplicationStatus) GormDataType() string   { return "application_status" }
func (ApplicationStatus) GormDBDataType() string { return "application_status" }
func (a ApplicationStatus) String() string       { return string(a) }
func AutoMigrateApplicationStatus(db *gorm.DB) error {
	if err := AutoMigrateEnumType(
		"application_status", db, StatusPending, StatusOffered, StatusUnknown, StatusRejected, StatusAccepted,
	); err != nil {
		return err
	}

	return nil
}
func (a *ApplicationStatus) UnmarshallJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &a); err != nil {
		return err
	}
	switch s {
	case "accepted":
		*a = StatusAccepted
	case "rejected":
		*a = StatusRejected
	case "pending":
		*a = StatusPending
	case "offered":
		*a = StatusOffered
	case "unknown":
		*a = StatusUnknown
	default:
		return fmt.Errorf("invalid application status type %s. Allowed: accepted, rejected, pending, offered, unknown", s)
	}
	return nil
}

type EventApplicationStatus string

const (
	EventDraft     EventApplicationStatus = "draft"
	EventOpen      EventApplicationStatus = "open"
	EventClosed    EventApplicationStatus = "closed"
	EventCancelled EventApplicationStatus = "cancelled"
	EventUnknown   EventApplicationStatus = "unknown"
)

func (EventApplicationStatus) GormDataType() string   { return "event_application_status" }
func (EventApplicationStatus) GormDBDataType() string { return "event_application_status" }
func (e EventApplicationStatus) String() string       { return string(e) }
func AutoMigrateEventApplicationStatus(db *gorm.DB) error {
	if err := AutoMigrateEnumType(
		"event_application_status", db, EventOpen, EventUnknown, EventClosed, EventCancelled,
	); err != nil {
		return err
	}
	return nil
}
func (e *EventApplicationStatus) UnmarshallJson(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	switch s {
	case "draft":
		*e = EventDraft
	case "open":
		*e = EventOpen
	case "closed":
		*e = EventClosed
	case "cancelled":
		*e = EventCancelled
	case "unknown":
		*e = EventUnknown
	default:
		return fmt.Errorf("invalid event application status %s. Allowed: draft, open, closed, cancelled, unknown", s)
	}
	return nil
}

type Application struct {
	Model
	Name             string
	Status           ApplicationStatus `gorm:"type:application_status;default:'unknown'" json:"application_status,default='unknown'"`
	Performer        Profile           `gorm:"foreignKey:PerformerID"`
	PerformerID      uuid.UUID         `json:"-" gorm:"performer_id,type:uuid"`
	EventRef         uuid.UUID         `gorm:"event_ref;type:uuid"`
	GoogleResponseID GoogleResponseID
}

type Event struct {
	Model
	Name         string
	Description  string
	Tags         []Tag         `gorm:"many2many:event_tags;"`
	Producer     Profile       `gorm:"foreignKey:ProducerID"`
	ProducerID   uuid.UUID     `json:"-" gorm:"producer_id;type:uuid"`
	Venue        *Profile      `gorm:"foreignKey:VenueID"`
	VenueID      *uuid.UUID    `json:"-" gorm:"venue_id;type:uuid"`
	Applications []Application `gorm:"foreignKey:EventRef"`
	GoogleForm   GoogleFormID
	Location     gormGIS.GeoPoint
	Status       EventApplicationStatus `json:"application_status,default='unknown'" gorm:"type:event_application_status;default:unknown"`
	Time         time.Time
	ApplyByTime  *time.Time `json:"apply_by_time,omitempty"`
	PayStructure string     `json:"pay_structure,omitempty"`
}

package entity

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"time"
)

type ApplicationStatus string

type Location struct {
	Latitude  float32
	Longitude float32
}

var LocationNil Location = Location{0.0, 0.0}

const (
	ACCEPTED ApplicationStatus = "accepted"
	REJECTED ApplicationStatus = "rejected"
	PENDING  ApplicationStatus = "pending"
	OFFERED  ApplicationStatus = "offered"
	UNKNOWN  ApplicationStatus = "unknown"
)

type Application struct {
	ID        uuid.UUID
	Name      string
	Status    ApplicationStatus
	Performer *Performer
	Event     *Event
	Act       *Act
}

func NewApplication(name string, status ApplicationStatus, performer *Performer, event *Event, act *Act) (Application, error) {
	if performer == nil {
		return Application{}, errors.New("performer required to create application")
	}
	if name == "" {
		return Application{}, fmt.Errorf("no name for application")
	}
	if event == nil {
		return Application{}, fmt.Errorf("event required for application")
	}
	return Application{Name: name, Status: status, Performer: performer, Event: event, Act: act}, nil
}

func ApplicationFromContext(ctx context.Context) (Application, error) {
	return Application{}, errors.New("not implemented error")
}

type GoogleForm map[string]interface{}
type EventApplicationStatus string

const (
	DRAFT     EventApplicationStatus = "draft"
	OPEN      EventApplicationStatus = "open"
	CLOSED    EventApplicationStatus = "closed"
	CANCELLED EventApplicationStatus = "cancelled"
)

type Event struct {
	ID              uuid.UUID
	ApplyByTime     *time.Time
	Name            string
	Description     string
	Tags            []Tag
	Producer        *Producer
	Venue           *Venue
	ApplicationForm GoogleForm
	Applications    []*Application
	ConfirmedActs   []*Act
	Status          EventApplicationStatus
	Location        Location
	EventTime       time.Time
	PayStructure    string
}

func NewEvent(
	name string, description string, applyByTime *time.Time,
	tags []Tag, producer *Producer, form ApplicationForm,
	applications []*Application, confirmedActs []*Act,
	location *Location, venue *Venue, eventTime time.Time,
) (Event, error) {
	if name == "" {
		return Event{}, fmt.Errorf("name required for event")
	}
	if tags == nil || len(tags) == 0 {
		return Event{}, fmt.Errorf("event must have at least one tag")
	}
	if producer == nil {
		return Event{}, fmt.Errorf("event must have a producer")
	}

	if (location == nil || *location == LocationNil) && (venue == nil) {
		return Event{}, errors.New("either location or venue must be supplied")
	}
	e := Event{
		Name:            name,
		ApplyByTime:     applyByTime,
		Description:     description,
		Tags:            tags,
		Producer:        producer,
		ApplicationForm: form,
		Applications:    applications,
		ConfirmedActs:   confirmedActs,
	}
	return e, nil
}

func EventFromContext(ctx context.Context) (Event, error) {
	return Event{}, errors.New("not implemented")
}

type Act struct {
	ID          uuid.UUID
	Name        string
	Tags        []Tag
	Description string
	Media       []Media
	Form        ApplicationForm
	Performer   *Performer
	Event       *Event
}

type Media interface {
	Get() (interface{}, error)
}

func NewAct(name string, tags []Tag, description string, media []Media, form ApplicationForm, performer *Performer) (Act, error) {
	if performer == nil {
		return Act{}, errors.New("act must have performer")
	}
	if tags == nil || len(tags) == 0 {
		return Act{}, fmt.Errorf("act must have at least one tag")
	}
	return Act{
		Name:        name,
		Tags:        tags,
		Description: description,
		Media:       media,
		Form:        form,
		Performer:   performer,
	}, nil
}

func ActFromContext(ctx context.Context) (Act, error) {
	return Act{}, errors.New("not implemented")
}

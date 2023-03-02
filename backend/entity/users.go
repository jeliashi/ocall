package entity

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// Many to many user id to performer/producer/venues
// you have to think about admin/poweruser

type Performer struct {
	ID              uuid.UUID
	Name            string
	UserID          uuid.UUID
	Applications    []*Application
	OfferedEvents   []*Event
	ConfirmedEvents []*Event
	Acts            []*Act
}
type Producer struct {
	ID        uuid.UUID
	Name      string
	UserID    uuid.UUID
	Events    []*Event
	SavedActs []*Act
}
type Venue struct {
	ID     uuid.UUID
	Name   string
	UserID uuid.UUID
	Events []*Event
}

func NewPerformer(name string, userID uuid.UUID) (Performer, error) {
	if name == "" {
		return Performer{}, fmt.Errorf("no name for performer")
	}
	p := &Performer{
		Name: name, UserID: userID,
	}
	return *p, nil
}

func PerformerFromContext(ctx context.Context) (Performer, error) {
	return Performer{}, errors.New("not implemented")
}

func NewProducer(name string, userID uuid.UUID) (Producer, error) {
	if name == "" {
		return Producer{}, fmt.Errorf("no name for producer")
	}
	p := &Producer{
		Name: name, UserID: userID,
	}
	return *p, nil
}

func ProducerFromContext(ctx context.Context) (Producer, error) {
	return Producer{}, errors.New("not implemented")
}

func NewVenue(name string, userID uuid.UUID) (Venue, error) {
	if name == "" {
		return Venue{}, fmt.Errorf("no name for venue")
	}
	v := &Venue{
		Name: name, UserID: userID,
	}
	return *v, nil
}

func VenueFromContext(ctx context.Context) (Venue, error) {
	return Venue{}, errors.New("not implemented")
}

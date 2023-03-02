package users

import (
	"context"
	"github.com/google/uuid"
	"ocall/app/entity"
)

type UseCase interface {
	CreatePerformer(ctx context.Context, name string, userID uuid.UUID) (entity.Performer, error)
	GetPerformer(ctx context.Context, id uuid.UUID) (entity.Performer, error)
	GetSavedPerformersByProducer(ctx context.Context, producer entity.Producer) ([]entity.Performer, error)
	GetPerformersByEvent(ctx context.Context, event entity.Event) ([]entity.Performer, error)
	UpdatePerformer(ctx context.Context, performer entity.Performer) error
	DeletePerformer(ctx context.Context, id uuid.UUID) error

	CreateProducer(ctx context.Context, name string, userID uuid.UUID) (entity.Producer, error)
	GetProducer(ctx context.Context, id uuid.UUID) (entity.Producer, error)
	GetProducerByEvent(ctx context.Context, event *entity.Event) (entity.Producer, error)
	UpdateProducer(ctx context.Context, performer entity.Producer) error
	DeleteProducer(ctx context.Context, id uuid.UUID) error

	CreateVenue(ctx context.Context, name string, userID uuid.UUID) (entity.Venue, error)
	GetVenue(ctx context.Context, id uuid.UUID) (entity.Venue, error)
	GetVenueByEvent(ctx context.Context, event entity.Event) (entity.Venue, error)
	UpdateVenue(ctx context.Context, performer entity.Venue) error
	DeleteVenue(ctx context.Context, id uuid.UUID) error

	ListProfiles(ctx context.Context, userID uuid.UUID) ([]entity.Performer, []entity.Producer, []entity.Venue, error)
}

type Repository interface {
	CreatePerformer(ctx context.Context, performer entity.Performer) (uuid.UUID, error)
	GetPerformerByID(ctx context.Context, id uuid.UUID) (entity.Performer, error)
	GetPerformersByProducerID(ctx context.Context, id uuid.UUID) ([]entity.Performer, error)
	GetPerformersByEventID(ctx context.Context, id uuid.UUID) ([]entity.Performer, error)
	UpdatePerformer(ctx context.Context, performer entity.Performer) error
	DeletePerformer(ctx context.Context, id uuid.UUID) error

	CreateProducer(ctx context.Context, producer entity.Producer) (uuid.UUID, error)
	GetProducerByID(ctx context.Context, id uuid.UUID) (entity.Producer, error)
	GetProducerByEventID(ctx context.Context, eventID uuid.UUID) (entity.Producer, error)
	UpdateProducer(ctx context.Context, performer entity.Producer) error
	DeleteProducer(ctx context.Context, id uuid.UUID) error

	CreateVenue(ctx context.Context, venue entity.Venue) (uuid.UUID, error)
	GetVenueByID(ctx context.Context, id uuid.UUID) (entity.Venue, error)
	GetVenueByEventID(ctx context.Context, eventID uuid.UUID) (entity.Venue, error)
	UpdateVenue(ctx context.Context, venue entity.Venue) error
	DeleteVenue(ctx context.Context, id uuid.UUID) error

	GetPerformersByUser(ctx context.Context, userID uuid.UUID) ([]entity.Performer, error)
	GetProducersByUser(ctx context.Context, userID uuid.UUID) ([]entity.Producer, error)
	GetVenuesByUser(ctx context.Context, userID uuid.UUID) ([]entity.Venue, error)
}

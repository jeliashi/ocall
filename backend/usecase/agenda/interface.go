package agenda

import (
	"context"
	"github.com/google/uuid"
	"ocall/app/entity"
)

type UseCase interface {
	CreateApplication(
		context context.Context, name string, status entity.ApplicationStatus,
		performer entity.Performer, event entity.Event, act entity.Act,
	) (entity.Application, error)
	GetApplication(ctx context.Context, id uuid.UUID) (entity.Application, error)
	GetApplicationsByEvent(ctx context.Context, event entity.Event) ([]entity.Application, error)
	GetApplicationsByPerformer(ctx context.Context, performer entity.Performer) ([]entity.Application, error)
	UpdateApplication(ctx context.Context, application entity.Application) error
	DeleteApplication(ctx context.Context, id uuid.UUID) error

	CreateEvent(
		ctx context.Context, name, description string,
		tags []entity.Tag, producer entity.Producer,
		form entity.ApplicationForm,
	) (entity.Event, error)
	GetEvent(ctx context.Context, id uuid.UUID) (entity.Event, error)
	AddPerformersToEvent(ctx context.Context, event entity.Event, performers []entity.Performer) error
	GetEventsByProducer(ctx context.Context, producer entity.Producer) ([]entity.Event, error)
	GetEventsByVenue(ctx context.Context, venue entity.Venue) ([]entity.Act, error)
	UpdateEvent(ctx context.Context, event entity.Event) error
	ChangeEventApplicationStatus(ctx context.Context, id uuid.UUID, status entity.EventApplicationStatus) error
	ModifyActListForEvent(ctx context.Context, event entity.Event, acts []entity.Act) error
	DeleteEvent(ctx context.Context, id uuid.UUID) error

	CreateAct(
		ctx context.Context, name, description string,
		tags []entity.Tag, media []entity.Media,
		form entity.ApplicationForm,
	) (entity.Act, error)
	GetActs(ctx context.Context, ids ...uuid.UUID) ([]entity.Act, error)
	GetActsByPerformer(ctx context.Context, performer entity.Performer) ([]entity.Act, error)
	GetConfirmedActsByEvent(ctx context.Context, event entity.Event) ([]entity.Act, error)
	GetSavedActByProducer(ctx context.Context, producer entity.Producer) ([]entity.Act, error)
	UpdateAct(ctx context.Context, act entity.Act) error
	DeleteAct(ctx context.Context, id uuid.UUID) error
}

type Repository interface {
	CreateApplication(ctx context.Context, application entity.Application) (uuid.UUID, error)
	GetApplicationByID(ctx context.Context, id uuid.UUID) (entity.Application, error)
	GetApplicationByEventID(ctx context.Context, id uuid.UUID) ([]entity.Application, error)
	GetApplicationsByPerformerID(ctx context.Context, id uuid.UUID) ([]entity.Application, error)
	UpdateApplication(ctx context.Context, application entity.Application) error
	DeleteApplication(ctx context.Context, id uuid.UUID) error

	CreateEvent(ctx context.Context, event entity.Event) (uuid.UUID, error)
	GetEventByID(ctx context.Context, id uuid.UUID) (entity.Event, error)
	AddApplicationsToEvent(ctx context.Context, event entity.Event, performers []entity.Performer) error
	GetEventsByProducerID(ctx context.Context, id uuid.UUID) ([]entity.Event, error)
	GetEventsByPerformerID(ctx context.Context, id uuid.UUID) ([]entity.Event, error)
	GetEventsByVenueID(ctx context.Context, id uuid.UUID) ([]entity.Event, error)
	UpdateEvent(ctx context.Context, event entity.Event) error
	ChangeEventApplicationStatus(ctx context.Context, id uuid.UUID, status entity.EventApplicationStatus) error
	ModifyActListForEvent(ctx context.Context, eventID uuid.UUID, acts []entity.Act) error
	DeleteEvent(ctx context.Context, id uuid.UUID) error

	CreateAct(ctx context.Context, act entity.Act) (uuid.UUID, error)
	GetActs(ctx context.Context, ids []uuid.UUID) ([]entity.Act, error)
	GetActsByPerformerID(ctx context.Context, performerID uuid.UUID) ([]entity.Act, error)
	GetConfirmedActsByEventID(ctx context.Context, eventID uuid.UUID) ([]entity.Act, error)
	GetSavedActsByProducerID(ctx context.Context, producerID uuid.UUID) ([]entity.Act, error)
	UpdateAct(ctx context.Context, act entity.Act) error
	DeleteAct(ctx context.Context, id uuid.UUID) error
}

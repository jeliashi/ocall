package agenda

import (
	"backend/models"
	"context"
	"github.com/google/uuid"
	"github.com/nferruzzi/gormGIS"
	"time"
)

type Repository interface {
	CreateEvent(ctx context.Context, event models.Event) (uuid.UUID, error)
	GetEvent(ctx context.Context, id uuid.UUID) (models.Event, error)
	UpdateEvent(ctx context.Context, event models.Event) (models.Event, error)
	DeleteEvent(ctx context.Context, id uuid.UUID) error

	CreateApplication(ctx context.Context, application models.Application) (uuid.UUID, error)
	GetApplication(ctx context.Context, id uuid.UUID) (models.Application, error)
	UpdateApplication(ctx context.Context, application models.Application) (models.Application, error)
	DeleteApplication(ctx context.Context, id uuid.UUID) error

	GetEventsByProducer(ctx context.Context, producerID uuid.UUID) ([]models.Event, error)
	GetApplicationsByEvent(ctx context.Context, eventID uuid.UUID) ([]models.Application, error)
	GetApplicationsByPerformer(ctx context.Context, performerID uuid.UUID) ([]models.Application, error)
	GetAllEvents(ctx context.Context, startTime time.Time, endTime time.Time, centerPoint gormGIS.GeoPoint, distanceKM float64) ([]models.Event, error)

	CreateTag(ctx context.Context, tag models.Tag) (uint, error)
	DeleteTag(ctx context.Context, tag models.Tag) error
}

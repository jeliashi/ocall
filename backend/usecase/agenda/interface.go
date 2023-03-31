package agenda

import (
	"backend/models"
	"context"
	"github.com/google/uuid"
	gormGIS "github.com/nferruzzi/gormgis"
	"time"
)

type Repository interface {
	CreateEvent(ctx context.Context, event models.Event) (uuid.UUID, error)
	GetEvent(ctx context.Context, id uuid.UUID) (models.Event, error)
	UpdateEvent(ctx context.Context, event models.Event) error
	DeleteEvent(ctx context.Context, id uuid.UUID) error

	CreateApplication(ctx context.Context, application models.Application) (uuid.UUID, error)
	GetApplication(ctx context.Context, id uuid.UUID) (models.Application, error)
	UpdateApplication(ctx context.Context, application models.Application) error
	DeleteApplication(ctx context.Context, id uuid.UUID) error

	GetEventsByProducer(ctx context.Context, producerID uuid.UUID) ([]models.Event, error)
	GetEventsByPerformer(ctx context.Context, performerID uuid.UUID) ([]models.Event, error)
	GetAllEvents(ctx context.Context, startTime time.Time, endTime time.Time, centerPoint gormGIS.GeoPoint, distanceKM float32) ([]models.Event, error)
	GetApplicationsByPerformer(ctx context.Context, performerID uuid.UUID) ([]models.Application, error)

	SaveApplication(ctx context.Context, applicationID uuid.UUID) error
	UnSaveApplication(ctx context.Context, applicationID uuid.UUID) error
	UpdateApplicationStatus(ctx context.Context, applicationID uuid.UUID, status models.ApplicationStatus) error
	UpdateEventApplicationStatus(ctx context.Context, eventID uuid.UUID, status models.EventApplicationStatus) error

	CreateTag(ctx context.Context, tag models.Tag) (uint, error)
	DeleteTag(ctx context.Context, tag models.Tag) error
}

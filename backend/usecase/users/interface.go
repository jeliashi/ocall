package users

import (
	"backend/models"
	"context"
	"github.com/google/uuid"
)

type Repository interface {
	CreateProfile(ctx context.Context, profile models.Profile) (uuid.UUID, error)
	GetProfileByID(ctx context.Context, id uuid.UUID) (models.Profile, error)
	GetPerformersByProducerID(ctx context.Context, producerID uuid.UUID) ([]models.Profile, error)
	GetPerformersByEventID(ctx context.Context, id uuid.UUID) ([]models.Profile, error)
	UpdateProfile(ctx context.Context, performer models.Profile) error
	DeleteProfile(ctx context.Context, id uuid.UUID) error
	GetProducerByEventID(ctx context.Context, eventID uuid.UUID) (models.Profile, error)
	GetVenueByEventID(ctx context.Context, eventID uuid.UUID) (models.Profile, error)

	GetProfilesByUser(ctx context.Context) ([]models.Profile, error)
	GetProfilePermissionLevel(ctx context.Context, performerID uuid.UUID) (models.Permission, error)
}

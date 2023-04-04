package users

import (
	"backend/models"
	"context"
	"github.com/google/uuid"
)

type Repository interface {
	CreateProfile(ctx context.Context, profile models.Profile) (uuid.UUID, error)
	GetProfileByID(ctx context.Context, id uuid.UUID) (models.Profile, error)
	UpdateProfile(ctx context.Context, profile models.Profile) (models.Profile, error)
	DeleteProfile(ctx context.Context, id uuid.UUID) error
	GetUsersByProfileId(ctx context.Context, id uuid.UUID) ([]models.UserID, error)
}

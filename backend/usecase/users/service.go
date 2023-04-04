package users

import (
	"backend/models"
	"context"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type Service struct {
	repo Repository
}

func NewService(repository Repository) Service {
	return Service{repository}
}

func (s *Service) CreateProfile(ctx context.Context, profile models.Profile) (uuid.UUID, error) {
	id, err := s.repo.CreateProfile(ctx, profile)
	if err != nil {
		return uuid.Nil, errors.Wrap(err, "db error")
	}
	return id, nil
}
func (s *Service) GetProfileByID(ctx context.Context, id uuid.UUID) (models.Profile, error) {
	performer, err := s.repo.GetProfileByID(ctx, id)
	if err != nil {
		return models.Profile{}, errors.Wrap(err, "db error")
	}
	return performer, nil
}
func (s *Service) UpdateProfile(ctx context.Context, profile models.Profile) (models.Profile, error) {
	if out, err := s.repo.UpdateProfile(ctx, profile); err != nil {
		return profile, errors.Wrap(err, "db error")
	} else {
		return out, nil
	}

}
func (s *Service) DeleteProfile(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.DeleteProfile(ctx, id); err != nil {
		return errors.Wrap(err, "db error")
	}
	return nil
}
func (s *Service) GetUsersByProfileId(ctx context.Context, id uuid.UUID) ([]models.UserID, error) {
	if users, err := s.repo.GetUsersByProfileId(ctx, id); err != nil {
		return nil, errors.Wrap(err, "error getting profiles from repo")
	} else {
		return users, nil
	}
}

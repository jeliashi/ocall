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
func (s *Service) GetPerformersByProducerID(ctx context.Context, producerID uuid.UUID) ([]models.Profile, error) {
	performers, err := s.repo.GetPerformersByProducerID(ctx, producerID)
	if err != nil {
		return nil, errors.Wrap(err, "db error")
	}
	return performers, nil
}
func (s *Service) GetPerformersByEventID(ctx context.Context, id uuid.UUID) ([]models.Profile, error) {
	performers, err := s.repo.GetPerformersByEventID(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "db error")
	}
	return performers, nil
}
func (s *Service) UpdateProfile(ctx context.Context, performer models.Profile) error {
	if err := s.repo.UpdateProfile(ctx, performer); err != nil {
		return errors.Wrap(err, "db error")
	}
	return nil
}
func (s *Service) DeleteProfile(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.DeleteProfile(ctx, id); err != nil {
		return errors.Wrap(err, "db error")
	}
	return nil
}
func (s *Service) GetProducerByEventID(ctx context.Context, eventID uuid.UUID) (models.Profile, error) {
	producer, err := s.repo.GetProducerByEventID(ctx, eventID)
	if err != nil {
		return models.Profile{}, errors.Wrap(err, "db error")
	}
	return producer, nil
}
func (s *Service) GetVenueByEventID(ctx context.Context, eventID uuid.UUID) (models.Profile, error) {
	venue, err := s.repo.GetVenueByEventID(ctx, eventID)
	if err != nil {
		return models.Profile{}, errors.Wrap(err, "db error")
	}
	return venue, nil
}
func (s *Service) GetProfilesByUser(ctx context.Context) ([]models.Profile, error) {
	performers, err := s.repo.GetProfilesByUser(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "db error")
	}
	return performers, nil
}
func (s *Service) GetProfilePermissionLevel(ctx context.Context, performerID uuid.UUID) (models.Permission, error) {
	permission, err := s.repo.GetProfilePermissionLevel(ctx, performerID)
	if err != nil {
		return models.PermissionUnknown, errors.Wrap(err, "db err")
	}
	return permission, nil
}

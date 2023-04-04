package agenda

import (
	"backend/models"
	"context"
	"github.com/google/uuid"
	"github.com/nferruzzi/gormGIS"
	"github.com/pkg/errors"
	"time"
)

type Service struct {
	repo Repository
}

func NewService(repository Repository) Service {
	return Service{repository}
}

func (s *Service) CreateEvent(ctx context.Context, event models.Event) (uuid.UUID, error) {
	id, err := s.repo.CreateEvent(ctx, event)
	if err != nil {
		return uuid.Nil, errors.Wrap(err, "db error")
	}
	return id, nil
}

func (s *Service) GetEvent(ctx context.Context, id uuid.UUID) (models.Event, error) {
	event, err := s.repo.GetEvent(ctx, id)
	if err != nil {
		return models.Event{}, errors.Wrap(err, "db error")
	}
	return event, nil
}
func (s *Service) UpdateEvent(ctx context.Context, event models.Event) (models.Event, error) {
	if out, err := s.repo.UpdateEvent(ctx, event); err != nil {
		return event, errors.Wrap(err, "db error")
	} else {
		return out, nil
	}

}
func (s *Service) DeleteEvent(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.DeleteEvent(ctx, id); err != nil {
		return errors.Wrap(err, "db error")
	}
	return nil
}
func (s *Service) CreateApplication(ctx context.Context, application models.Application) (uuid.UUID, error) {
	id, err := s.repo.CreateApplication(ctx, application)
	if err != nil {
		return uuid.Nil, errors.Wrap(err, "db error")
	}
	return id, nil
}
func (s *Service) GetApplication(ctx context.Context, id uuid.UUID) (models.Application, error) {
	application, err := s.repo.GetApplication(ctx, id)
	if err != nil {
		return models.Application{}, errors.Wrap(err, "db error")
	}
	return application, nil
}
func (s *Service) UpdateApplication(ctx context.Context, application models.Application) (models.Application, error) {
	if out, err := s.repo.UpdateApplication(ctx, application); err != nil {
		return application, errors.Wrap(err, "db error")
	} else {
		return out, nil
	}
}
func (s *Service) DeleteApplication(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.DeleteApplication(ctx, id); err != nil {
		return errors.Wrap(err, "db error")
	}
	return nil
}
func (s *Service) GetEventsByProducer(ctx context.Context, producerID uuid.UUID) ([]models.Event, error) {
	events, err := s.repo.GetEventsByProducer(ctx, producerID)
	if err != nil {
		return nil, errors.Wrap(err, "db error")
	}
	return events, nil
}

func (s *Service) GetApplicationsByEvent(ctx context.Context, eventID uuid.UUID) ([]models.Application, error) {
	applications, err := s.repo.GetApplicationsByEvent(ctx, eventID)
	if err != nil {
		return nil, errors.Wrap(err, "db error")
	}
	return applications, nil
}

func (s *Service) GetApplicationsByPerformer(ctx context.Context, performerID uuid.UUID) ([]models.Application, error) {
	applications, err := s.repo.GetApplicationsByPerformer(ctx, performerID)
	if err != nil {
		return nil, errors.Wrap(err, "db error")
	}
	return applications, nil
}
func (s *Service) GetAllEvents(
	ctx context.Context, startTime time.Time, endTime time.Time,
	centerPoint gormGIS.GeoPoint, distanceKM float64,
) ([]models.Event, error) {
	events, err := s.repo.GetAllEvents(ctx, startTime, endTime, centerPoint, distanceKM)
	if err != nil {
		return nil, errors.Wrap(err, "db error")
	}
	return events, nil
}

func (s *Service) CreateTag(ctx context.Context, tag models.Tag) (uint, error) {
	id, err := s.repo.CreateTag(ctx, tag)
	if err != nil {
		return 0, errors.Wrap(err, "db error")
	}
	return id, nil
}
func (s *Service) DeleteTag(ctx context.Context, tag models.Tag) error {
	if err := s.repo.DeleteTag(ctx, tag); err != nil {
		return errors.Wrap(err, "db error")
	}
	return nil
}

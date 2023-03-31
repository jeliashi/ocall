package agenda

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	gormGIS "github.com/nferruzzi/gormgis"
	"github.com/pkg/errors"
	"ocall/backend/models"
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
func (s *Service) UpdateEvent(ctx context.Context, event models.Event) error {
	if err := s.repo.UpdateEvent(ctx, event); err != nil {
		return errors.Wrap(err, "db error")
	}
	return nil
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
func (s *Service) UpdateApplication(ctx context.Context, application models.Application) error {
	if err := s.repo.UpdateApplication(ctx, application); err != nil {
		return errors.Wrap(err, "db error")
	}
	return nil
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
func (s *Service) GetEventsByPerformer(ctx context.Context, performerID uuid.UUID) ([]models.Event, error) {
	events, err := s.repo.GetEventsByPerformer(ctx, performerID)
	if err != nil {
		return nil, errors.Wrap(err, "db error")
	}
	return events, nil
}
func (s *Service) GetAllEvents(
	ctx context.Context, startTime time.Time, endTime time.Time,
	centerPoint gormGIS.GeoPoint, distanceKM float32,
) ([]models.Event, error) {
	events, err := s.repo.GetAllEvents(ctx, startTime, endTime, centerPoint, distanceKM)
	if err != nil {
		return nil, errors.Wrap(err, "db error")
	}
	return events, nil
}
func (s *Service) GetApplicationsByPerfomer(ctx context.Context, performerID uuid.UUID) ([]models.Application, error) {
	applications, err := s.repo.GetApplicationsByPerformer(ctx, performerID)
	if err != nil {
		return nil, errors.Wrap(err, "db error")
	}
	return applications, nil
}
func (s *Service) GetApplicationsByProducer(ctx context.Context, producerID uuid.UUID) ([]models.Application, error) {
	applications, err := s.repo.GetApplicationsByProducer(ctx, producerID)
	if err != nil {
		return nil, errors.Wrap(err, "db error")
	}
	return applications, nil
}
func (s *Service) SaveApplication(ctx context.Context, applicationID uuid.UUID) error {
	application, err := s.repo.GetApplication(ctx, applicationID)
	if err != nil {
		return errors.Wrapf(err, "unable to retrieve application id %s", applicationID)
	}
	if application.Saved {
		return errors.New("application already saved")
	}

	if err := s.repo.SaveApplication(ctx, applicationID); err != nil {
		return errors.Wrap(err, "unable to save application")
	}
	return nil
}
func (s *Service) UnSaveApplication(ctx context.Context, applicationID uuid.UUID) error {
	application, err := s.repo.GetApplication(ctx, applicationID)
	if err != nil {
		return errors.Wrapf(err, "unable to retrieve application id %s", applicationID)
	}
	if !application.Saved {
		return errors.New("application already saved")
	}

	if err := s.repo.UnSaveApplication(ctx, applicationID); err != nil {
		return errors.Wrap(err, "unable to un-save application")
	}
	return nil
}
func (s *Service) UpdateApplicationStatus(
	ctx context.Context, applicationID uuid.UUID, status models.ApplicationStatus,
) error {
	application, err := s.repo.GetApplication(ctx, applicationID)
	if err != nil {
		return errors.Wrapf(err, "unable to retrieve application %s", applicationID)
	}
	if application.Status == status {
		return fmt.Errorf("application already has status %s", string(status))
	}
	if err := s.repo.UpdateApplicationStatus(ctx, applicationID, status); err != nil {
		return errors.Wrap(err, "db error")
	}
	return nil
}
func (s *Service) UpdateEventApplicationStatus(
	ctx context.Context, eventID uuid.UUID, status models.EventApplicationStatus,
) error {
	event, err := s.repo.GetEvent(ctx, eventID)
	if err != nil {
		return errors.Wrapf(err, "unable to retrieve event %s", eventID)
	}
	if event.Status == status {
		return fmt.Errorf("event already has status %s", string(status))
	}
	if err := s.repo.UpdateEventApplicationStatus(ctx, eventID, status); err != nil {
		return errors.Wrap(err, "db error")
	}
	return nil
}

func (s *Service) CreateTag(ctx context.Context, tag models.Tag) error {
	if err := s.repo.CreateTag(ctx, tag); err != nil {
		return errors.Wrap(err, "db error")
	}
	return nil
}
func (s *Service) DeleteTag(ctx context.Context, tag models.Tag) error {
	if err := s.repo.DeleteTag(ctx, tag); err != nil {
		return errors.Wrap(err, "db error")
	}
	return nil
}

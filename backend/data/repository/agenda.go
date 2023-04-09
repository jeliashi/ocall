package repository

import (
	"backend/models"
	"context"
	"github.com/google/uuid"
	"github.com/nferruzzi/gormGIS"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"time"
)

type AgendaRepo struct {
	orm *gorm.DB
}

func NewAgendaRepo(db *gorm.DB) AgendaRepo {
	return AgendaRepo{orm: db}
}

func (r *AgendaRepo) CreateEvent(ctx context.Context, event models.Event) (uuid.UUID, error) {
	if err := r.orm.WithContext(ctx).Create(&event).Error; err != nil {
		return uuid.Nil, errors.Wrap(err, "gorm create error")
	}
	return event.ID, nil
}
func (r *AgendaRepo) GetEvent(ctx context.Context, id uuid.UUID) (models.Event, error) {
	var event models.Event
	if err := r.orm.WithContext(ctx).First(&event, id).Error; err != nil {
		return event, errors.Wrap(err, "gorm first error")
	}
	return event, nil
}
func (r *AgendaRepo) UpdateEvent(ctx context.Context, event models.Event) (models.Event, error) {
	if err := r.orm.WithContext(ctx).Save(&event).Error; err != nil {
		return event, errors.Wrap(err, "gorm save error")
	}
	return event, nil
}
func (r *AgendaRepo) DeleteEvent(ctx context.Context, id uuid.UUID) error {
	if err := r.orm.WithContext(ctx).Delete(&models.Event{}, id).Error; err != nil {
		return errors.Wrap(err, "gorm delete error")
	}
	return nil
}

func (r *AgendaRepo) CreateApplication(ctx context.Context, application models.Application) (uuid.UUID, error) {
	if err := r.orm.WithContext(ctx).Create(&application).Error; err != nil {
		return uuid.Nil, errors.Wrap(err, "gorm create error")
	}
	return application.ID, nil
}
func (r *AgendaRepo) GetApplication(ctx context.Context, id uuid.UUID) (models.Application, error) {
	var application models.Application
	if err := r.orm.WithContext(ctx).First(&application, id).Error; err != nil {
		return application, errors.Wrap(err, "gorm first error")
	}
	return application, nil
}
func (r *AgendaRepo) UpdateApplication(ctx context.Context, application models.Application) (models.Application, error) {
	if err := r.orm.WithContext(ctx).Save(&application).Error; err != nil {
		return application, errors.Wrap(err, "gorm save error")
	}
	return application, nil
}
func (r *AgendaRepo) DeleteApplication(ctx context.Context, id uuid.UUID) error {
	if err := r.orm.WithContext(ctx).Delete(&models.Application{}, id).Error; err != nil {
		return errors.Wrap(err, "gorm delete error")
	}
	return nil
}

func (r *AgendaRepo) GetEventsByProducer(ctx context.Context, producerID uuid.UUID) ([]models.Event, error) {
	var eventPointers []*models.Event
	if err := r.orm.WithContext(ctx).Where("producer_id = ?", producerID).Find(&eventPointers).Error; err != nil {
		return nil, errors.Wrap(err, "gorm find error")
	}
	events := make([]models.Event, len(eventPointers))
	for i, event := range eventPointers {
		events[i] = *event
	}
	return events, nil
}

func (r *AgendaRepo) GetApplicationsByEvent(ctx context.Context, eventID uuid.UUID) ([]models.Application, error) {
	var event models.Event
	if err := r.orm.WithContext(ctx).Preload("Applications").First(&event, eventID).Error; err != nil {
		return nil, errors.Wrap(err, "gorm first error")
	}
	return event.Applications, nil
}
func (r *AgendaRepo) GetApplicationsByPerformer(ctx context.Context, performerID uuid.UUID) ([]models.Application, error) {
	var applicationPointers []*models.Application
	if err := r.orm.WithContext(ctx).Where("performer_id = ?", performerID).Find(&applicationPointers).Error; err != nil {
		return nil, errors.Wrap(err, "gorm find error")
	}
	applications := make([]models.Application, len(applicationPointers))
	for j, app := range applicationPointers {
		applications[j] = *app
	}
	return applications, nil
}

func (r *AgendaRepo) GetAllEvents(
	ctx context.Context, startTime time.Time, endTime time.Time, centerPoint gormGIS.GeoPoint, distanceKM float64,
) ([]models.Event, error) {
	var eventPointers []*models.Event
	if err := r.orm.WithContext(ctx).Where("time >= ?", startTime).
		Where("time <= ?", endTime).
		Where("ST_Distance_Sphere(location, ?) <= ?", centerPoint, 1000.0*distanceKM).
		Find(&eventPointers).Error; err != nil {
		return nil, errors.Wrap(err, "gorm find error")
	}
	events := make([]models.Event, len(eventPointers))
	for j, e := range eventPointers {
		events[j] = *e
	}
	return events, nil
}

func (r *AgendaRepo) CreateTag(ctx context.Context, tag models.Tag) (uint, error) {
	if err := r.orm.WithContext(ctx).Create(&tag).Error; err != nil {
		return 0, errors.Wrap(err, "gorm create error")
	}
	return tag.ID, nil
}
func (r *AgendaRepo) DeleteTag(ctx context.Context, tag models.Tag) error {
	if err := r.orm.WithContext(ctx).Delete(&tag).Error; err != nil {
		return errors.Wrap(err, "gorm delete error")
	}
	return nil
}

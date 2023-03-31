package repository

import (
	"backend/models"
	"context"
	"github.com/google/uuid"
	gormGIS "github.com/nferruzzi/gormgis"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"strconv"
	"time"
)

type AgendaRepo struct {
	orm *gorm.DB
}

func NewAgendaRepo(db *gorm.DB) AgendaRepo {
	return AgendaRepo{orm: db}
}

func (r *AgendaRepo) checkTagsExist(ctx context.Context, tags []models.Tag) error {
	var dbTags []models.Tag
	filter := make([]string, len(tags))
	for j, tag := range tags {
		filter[j] = tag.Name
	}
	result := r.orm.WithContext(ctx).Where("name IN ?", filter).Find(&dbTags)
	if err := models.ParseGormErrors(result.Error, "tags", ""); err != nil {
		return errors.Wrap(err, "gorm error")
	}
	if len(dbTags) == len(tags) {
		return nil
	}
	missingTagNames := make([]string, 0)
	var contained bool
	for _, tag := range tags {
		contained = false
		for _, remoteTag := range dbTags {
			if remoteTag.Name == tag.Name {
				contained = true
			}
		}
		if !contained {
			missingTagNames = append(missingTagNames, tag.Name)
		}
	}
	return errors.Errorf("tags %s not recorded", missingTagNames)
}

func (r *AgendaRepo) getTagsByName(ctx context.Context, tags []models.Tag) ([]models.Tag, error) {
	var dbTags []models.Tag
	filter := make([]string, len(tags))
	for j, tag := range tags {
		filter[j] = tag.Name
	}
	result := r.orm.WithContext(ctx).Where("name IN ?", filter).Find(&dbTags)
	if err := models.ParseGormErrors(result.Error, "tags", ""); err != nil {
		return nil, errors.Wrap(err, "gorm error")
	}
	if len(dbTags) == len(tags) {
		return dbTags, nil
	}
	missingTagNames := make([]string, 0)
	var contained bool
	for _, tag := range tags {
		contained = false
		for _, remoteTag := range dbTags {
			if remoteTag.Name == tag.Name {
				contained = true
			}
		}
		if !contained {
			missingTagNames = append(missingTagNames, tag.Name)
		}
	}
	return nil, errors.Errorf("tags %s not recorded", missingTagNames)
}

func (r *AgendaRepo) CreateEvent(ctx context.Context, event models.Event) (uuid.UUID, error) {
	if err := r.checkTagsExist(ctx, event.Tags); err != nil {
		return uuid.Nil, err
	}
	result := r.orm.WithContext(ctx).Create(event)
	if result.Error != nil {
		return uuid.Nil, errors.Wrap(result.Error, "gorm error")
	}
	return event.ID, nil
}
func (r *AgendaRepo) GetEvent(ctx context.Context, id uuid.UUID) (models.Event, error) {
	var event models.Event

	//r.orm.WithContext(ctx).Model(&models.Producer{}).Preload("user_ids").Find()
	result := r.orm.WithContext(ctx).First(&event, id)
	if err := models.ParseGormErrors(result.Error, "events", id.String()); err != nil {
		return event, errors.Wrap(err, "gorm error")
	}
	return event, nil
}

func (r *AgendaRepo) UpdateEvent(ctx context.Context, event models.Event) error {
	updates := models.Event{
		Name:         event.Name,
		Description:  event.Description,
		Tags:         event.Tags,
		Venue:        event.Venue,
		Applications: event.Applications,
		GoogleForm:   event.GoogleForm,
		Location:     event.Location,
		Status:       event.Status,
		Time:         event.Time,
		ApplyByTime:  event.ApplyByTime,
		PayStructure: event.PayStructure,
	}
	result := r.orm.WithContext(ctx).Model(&event).Updates(updates)
	if err := models.ParseGormErrors(result.Error, "events", event.ID.String()); err != nil {
		return err
	}

	return nil
}
func (r *AgendaRepo) DeleteEvent(ctx context.Context, id uuid.UUID) error {
	var event models.Event
	result := r.orm.WithContext(ctx).Delete(&event, id)
	if err := models.ParseGormErrors(result.Error, "events", id.String()); err != nil {
		return err
	}
	return nil
}
func (r *AgendaRepo) CreateApplication(ctx context.Context, application models.Application) (uuid.UUID, error) {
	result := r.orm.WithContext(ctx).Create(application)
	if result.Error != nil {
		return uuid.Nil, errors.Wrap(result.Error, "gorm error")
	}
	return application.ID, nil
}
func (r *AgendaRepo) GetApplication(ctx context.Context, id uuid.UUID) (models.Application, error) {
	var application models.Application
	result := r.orm.WithContext(ctx).First(&application, id)
	if err := models.ParseGormErrors(result.Error, "applications", id.String()); err != nil {
		return application, err
	}
	return application, nil
}
func (r *AgendaRepo) UpdateApplication(ctx context.Context, application models.Application) error {
	updates := models.Application{
		Name:             application.Name,
		Status:           application.Status,
		GoogleResponseID: application.GoogleResponseID,
	}
	result := r.orm.WithContext(ctx).Model(&application).Updates(updates)
	if err := models.ParseGormErrors(result.Error, "applications", application.ID.String()); err != nil {
		return err
	}

	return nil
}
func (r *AgendaRepo) DeleteApplication(ctx context.Context, id uuid.UUID) error {
	var application models.Application
	result := r.orm.WithContext(ctx).Delete(&application, id)
	if err := models.ParseGormErrors(result.Error, "applications", id.String()); err != nil {
		return err
	}
	return nil
}
func (r *AgendaRepo) GetEventsByProducer(ctx context.Context, producerID uuid.UUID) ([]models.Event, error) {
	var events []models.Event
	result := r.orm.WithContext(ctx).Find(events, "producerRef = (?)", producerID)
	if err := models.ParseGormErrors(result.Error, "events", producerID.String()); err != nil {
		return nil, err
	}
	return events, nil
}
func (r *AgendaRepo) GetEventsByPerformer(ctx context.Context, performerID uuid.UUID) ([]models.Event, error) {
	var events []models.Event

	result := r.orm.WithContext(ctx).
		Preload("Applications").
		Where("application.performer_ref = ?", performerID.String()).
		Find(&events)
	if err := models.ParseGormErrors(result.Error, "events", performerID.String()); err != nil {
		return nil, err
	}

	return events, nil
}
func (r *AgendaRepo) GetAllEvents(
	ctx context.Context, startTime time.Time, endTime time.Time, centerPoint gormGIS.GeoPoint, distanceKM float32,
) ([]models.Event, error) {
	var events []models.Event
	query := r.orm.WithContext(ctx)
	if !startTime.Equal(time.Time{}) {
		query = query.Where("start_time >= ?", startTime)
	}
	if !endTime.Equal(time.Time{}) {
		query = query.Where("end_time <= ?", endTime)
	}
	if centerPoint.Lat != 0.0 && centerPoint.Lng != 0.0 && distanceKM > 0 {
		query = query.Where("ST_DWithin(?, ?, true)", centerPoint, 1000.0*distanceKM)
	}
	result := query.Find(&events)
	if err := models.ParseGormErrors(result.Error, "events", uuid.Nil.String()); err != nil {
		return nil, err
	}
	return events, nil
}
func (r *AgendaRepo) GetApplicationsByPerformer(
	ctx context.Context, performerID uuid.UUID,
) ([]models.Application, error) {
	var applications []models.Application
	result := r.orm.WithContext(ctx).Where("performer_ref = ?", performerID.String()).Find(&applications)
	if err := models.ParseGormErrors(result.Error, "performers", performerID.String()); err != nil {
		return nil, err
	}
	return applications, nil
}
func (r *AgendaRepo) GetApplicationsByEvent(ctx context.Context, eventID uuid.UUID) ([]models.Application, error) {
	var event models.Event
	result := r.orm.WithContext(ctx).Joins("Applications").Find(&event, eventID)
	if err := models.ParseGormErrors(result.Error, "events", eventID.String()); err != nil {
		return nil, err
	}
	return event.Applications, nil
}
func (r *AgendaRepo) SaveApplication(ctx context.Context, applicationID uuid.UUID) error {
	result := r.orm.WithContext(ctx).
		Model(&models.Application{}).
		Where("id = ?", applicationID).
		UpdateColumn("saved", true)
	if err := models.ParseGormErrors(result.Error, "applications", applicationID.String()); err != nil {
		return err
	}
	return nil
}
func (r *AgendaRepo) UnSaveApplication(ctx context.Context, applicationID uuid.UUID) error {
	result := r.orm.WithContext(ctx).
		Model(&models.Application{}).
		Where("id = ?", applicationID).
		UpdateColumn("saved", false)
	if err := models.ParseGormErrors(result.Error, "applications", applicationID.String()); err != nil {
		return err
	}
	return nil
}
func (r *AgendaRepo) UpdateApplicationStatus(
	ctx context.Context, applicationID uuid.UUID, status models.ApplicationStatus,
) error {
	result := r.orm.WithContext(ctx).
		Model(&models.Application{}).
		Where("id = ?", applicationID).
		UpdateColumn("status", status)
	if err := models.ParseGormErrors(result.Error, "applications", applicationID.String()); err != nil {
		return err
	}
	return nil
}
func (r *AgendaRepo) UpdateEventApplicationStatus(
	ctx context.Context, eventID uuid.UUID, status models.EventApplicationStatus,
) error {
	result := r.orm.WithContext(ctx).
		Model(&models.Event{}).
		Where("id = ?", eventID).
		UpdateColumn("status", status)
	if err := models.ParseGormErrors(result.Error, "events", eventID.String()); err != nil {
		return err
	}
	return nil
}
func (r *AgendaRepo) GetTags(ctx context.Context) ([]models.Tag, error) {
	var tags []models.Tag
	result := r.orm.WithContext(ctx).Find(&tags)
	if err := models.ParseGormErrors(result.Error, "tags", uuid.Nil.String()); err != nil {
		return nil, err
	}
	return tags, nil
}
func (r *AgendaRepo) CreateTag(ctx context.Context, tag models.Tag) (uint, error) {
	result := r.orm.WithContext(ctx).Create(&tag)
	if err := models.ParseGormErrors(result.Error, "tags", uuid.Nil.String()); err != nil {
		return 0, err
	}
	return tag.ID, nil
}
func (r *AgendaRepo) DeleteTag(ctx context.Context, tag models.Tag) error {
	result := r.orm.WithContext(ctx).Delete(&tag)
	if err := models.ParseGormErrors(result.Error, "tags", strconv.Itoa(int(tag.ID))); err != nil {
		return err
	}
	return nil
}

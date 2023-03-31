package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"ocall/backend/models"
)

type UserRepo struct {
	orm *gorm.DB
}

func NewUserRepo(db *gorm.DB) UserRepo {
	return UserRepo{orm: db}
}

func (r *UserRepo) CreateProfile(ctx context.Context, performer models.Profile) (uuid.UUID, error) {
	result := r.orm.WithContext(ctx).Create(&performer)
	if err := models.ParseGormErrors(result.Error, "performers", performer.ID.String()); err != nil {
		return uuid.Nil, err
	}
	return performer.ID, nil
}

func (r *UserRepo) GetProfileByID(ctx context.Context, id uuid.UUID) (models.Profile, error) {
	var profile models.Profile
	result := r.orm.WithContext(ctx).First(&profile, id)
	if err := models.ParseGormErrors(result.Error, "profiles", id.String()); err != nil {
		return models.Profile{}, err
	}
	return profile, nil
}

func (r *UserRepo) GetPerformersByProducerID(ctx context.Context, producerID uuid.UUID) ([]models.Profile, error) {
	var events []models.Event

	result := r.orm.WithContext(ctx).
		Preload("Applications").
		Preload("Applications.Performer").
		Where("producer_id = ?", producerID.String()).
		Find(&events)
	if err := models.ParseGormErrors(result.Error, "events", producerID.String()); err != nil {
		return nil, err
	}
	performers := make([]models.Profile, 0)
	for _, e := range events {
		for _, a := range e.Applications {
			if a.Status != models.STATUS_REJECTED {
				performers = append(performers, a.Performer)
			}
		}
	}
	return performers, nil
}

func (r *UserRepo) GetPerformersByEventID(ctx context.Context, id uuid.UUID) ([]models.Profile, error) {
	var applications []models.Application
	result := r.orm.WithContext(ctx).Preload("Performer").Where("event_id = ?", id.String()).Find(applications)
	if err := models.ParseGormErrors(result.Error, "applications", id.String()); err != nil {
		return nil, err
	}
	performers := make([]models.Profile, 0)
	for _, a := range applications {
		performers = append(performers, a.Performer)
	}
	return performers, nil
}

func (r *UserRepo) UpdateProfile(ctx context.Context, performer models.Profile) error {
	updates := models.Profile{
		Name:        performer.Name,
		ProfileType: "",
		Location:    performer.Location,
		UserIDs:     performer.UserIDs,
	}
	result := r.orm.Where(ctx).Model(&performer).Updates(updates)
	if err := models.ParseGormErrors(result.Error, "profiles", performer.ID.String()); err != nil {
		return err
	}
	return nil
}

func (r *UserRepo) DeleteProfile(ctx context.Context, id uuid.UUID) error {
	result := r.orm.WithContext(ctx).Delete(&models.Profile{Model: models.Model{ID: id}})
	if err := models.ParseGormErrors(result.Error, "profiles", id.String()); err != nil {
		return err
	}
	return nil
}

func (r *UserRepo) GetProducerByEventID(ctx context.Context, eventID uuid.UUID) (models.Profile, error) {
	var event models.Event
	result := r.orm.WithContext(ctx).Preload("Producer").First(&event, eventID)
	if err := models.ParseGormErrors(result.Error, "events", eventID.String()); err != nil {
		return models.Profile{}, err
	}
	return event.Producer, nil
}

func (r *UserRepo) GetVenueByEventID(ctx context.Context, eventID uuid.UUID) (models.Profile, error) {
	var event models.Event
	result := r.orm.WithContext(ctx).Preload("Venue").First(&event, eventID)
	if err := models.ParseGormErrors(result.Error, "events", eventID.String()); err != nil {
		return models.Profile{}, err
	}
	return *event.Venue, nil
}

func (r *UserRepo) GetProfilesByUser(ctx context.Context) ([]models.Profile, error) {
	var profiles []models.Profile
	userId := ctx.Value(models.FirebaseContextKey).(uuid.UUID)
	if userId == uuid.Nil {
		return nil, errors.New("no firebase user in context")
	}
	result := r.orm.WithContext(ctx).Preload("user_ids").Where("user_id.firebase_id = ?", userId.String()).Find(&profiles)
	if err := models.ParseGormErrors(result.Error, "user_ids", userId.String()); err != nil {
		return nil, err
	}
	return profiles, nil
}

func (r *UserRepo) GetProfilePermissionLevel(ctx context.Context, profileID uuid.UUID) (models.Permission, error) {
	var permissions models.Permission
	userID := models.GetFirebaseIDFromContext(ctx)
	result := r.orm.WithContext(ctx).
		Where("firebase = ? and profile_id = ?", userID, profileID).
		Pluck("permissions", &permissions)
	if err := models.ParseGormErrors(result.Error, "user_ids", profileID.String()); err != nil {
		return models.PermissionUnknown, err
	}
	return permissions, nil
}

func (r *UserRepo) GetApplicationsByProducerID(ctx context.Context, producerID uuid.UUID) ([]models.Application, error) {
	return nil, errors.New("unimplemented error")
}

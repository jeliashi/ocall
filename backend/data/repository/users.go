package repository

import (
	"backend/models"
	"context"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type UserRepo struct {
	orm *gorm.DB
}

func NewUserRepo(db *gorm.DB) UserRepo {
	return UserRepo{orm: db}
}

func (r *UserRepo) CreateProfile(ctx context.Context, profile models.Profile) (uuid.UUID, error) {
	result := r.orm.WithContext(ctx).Create(&profile)
	if result.Error != nil {
		return uuid.Nil, errors.Wrap(result.Error, "gorm create error")
	}
	return profile.ID, nil
}
func (r *UserRepo) GetProfileByID(ctx context.Context, id uuid.UUID) (models.Profile, error) {
	var profile models.Profile
	if err := r.orm.WithContext(ctx).First(&profile, id).Error; err != nil {
		return profile, errors.Wrap(err, "gorm first error")
	}
	return profile, nil
}
func (r *UserRepo) UpdateProfile(ctx context.Context, profile models.Profile) (models.Profile, error) {
	if err := r.orm.WithContext(ctx).Save(&profile).Error; err != nil {
		return profile, errors.Wrap(err, "gorm save error")
	}
	return profile, nil
}
func (r *UserRepo) DeleteProfile(ctx context.Context, id uuid.UUID) error {
	if err := r.orm.WithContext(ctx).Delete(&models.Profile{}, id).Error; err != nil {
		return errors.Wrap(err, "gorm delete error")
	}
	return nil
}
func (r *UserRepo) GetUsersByProfileId(ctx context.Context, id uuid.UUID) ([]models.UserID, error) {
	var profile models.Profile
	if err := r.orm.WithContext(ctx).Preload("UserID").Where("id = ?", id).First(&profile).Error; err != nil {
		return nil, errors.Wrap(err, "gorm find error")
	}
	return profile.UserIDs, nil
}

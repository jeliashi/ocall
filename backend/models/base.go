package models

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Permission string

const (
	Admin             Permission = "admin"
	Restricted                   = "restricted"
	PermissionUnknown            = "unknown"
)

type Model struct {
	ID        uuid.UUID      `json:"id" gorm:"primaryKey,unique"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

type NotFoundErr struct {
	db string
	id string
}

func (e *NotFoundErr) Error() string {
	return fmt.Sprintf("record %s not found in %s", e.id, e.db)
}

func ParseGormErrors(err error, db, id string) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &NotFoundErr{db, id}
	}
	return err
}

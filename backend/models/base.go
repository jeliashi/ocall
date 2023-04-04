package models

import (
	"encoding/json"
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

func (p *Permission) UnmarshallJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	switch s {
	case "admin":
		*p = Admin
	case "restricted":
		*p = Restricted
	default:
		*p = PermissionUnknown
	}
	return nil
}

type Model struct {
	ID        uuid.UUID      `json:"id" gorm:"primaryKey,unique"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

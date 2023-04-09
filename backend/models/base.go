package models

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"strings"
	"time"
)

type enumType interface {
	String() string
}

func AutoMigrateEnumType(name string, db *gorm.DB, elements ...enumType) error {
	enums := make([]string, len(elements))
	for j, element := range elements {
		enums[j] = element.String()
	}
	if result := db.Exec(fmt.Sprintf("SELECT 1 FROM pg_type WHERE typname = '%s'", name)); result.RowsAffected == 0 {
		if err := db.Exec(fmt.Sprintf(
			"CREATE TYPE %s AS ENUM ('%s')",
			name,
			strings.Join(enums, "', '"),
		)).Error; err != nil {
			return err
		}
	} else if result.Error != nil {
		return result.Error
	}
	return nil
}

type Permission string

const (
	Admin             Permission = "admin"
	Restricted        Permission = "restricted"
	PermissionUnknown Permission = "unknown"
)

func (p Permission) String() string { return string(p) }

func (Permission) GormDataType() string   { return "permission" }
func (Permission) GormDBDataType() string { return "permission" }
func AutoMigratePermission(db *gorm.DB) error {
	if err := AutoMigrateEnumType("permission", db, Admin, Restricted, PermissionUnknown); err != nil {
		return err
	}
	return nil
}

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
	ID        uuid.UUID      `json:"id,omitempty" gorm:"primaryKey,type:uuid,default:uuid_generate_v4()" swaggerignore:"true"`
	CreatedAt time.Time      `json:"created_at" swaggerignore:"true"`
	UpdatedAt time.Time      `json:"updated_at" swaggerignore:"true"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at" swaggerignore:"true"`
}

func (m *Model) BeforeCreate(tx *gorm.DB) error {
	m.ID = uuid.New()
	return nil
}

package models

import (
	"context"
	"github.com/google/uuid"
	gormGIS "github.com/nferruzzi/gormgis"
	"strings"
)

const FirebaseContextKey string = "firebase_context_key"

func GetFirebaseIDFromContext(ctx context.Context) string {
	return ctx.Value(FirebaseContextKey).(string)
}

type ProfileType string

const (
	ProducerType  ProfileType = "producer"
	PerformerType             = "performer"
	VenueType                 = "venue"
)

type UserID struct {
	Model
	FirebaseId  string
	Permissions Permission
	ProfileId   uuid.UUID
}

type Profile struct {
	Model
	Name        string `json:"name,nonempty" gorm:"notnull"`
	ProfileType `json:"type"`
	Location    *gormGIS.GeoPoint
	UserIDs     []UserID `json:"user_ids" gorm:"foreignKey:ProfileId"`
}

func ParseProfile(s string) (ProfileType, bool) {
	converter := map[string]ProfileType{"producer": ProducerType, "performer": PerformerType, "venue": VenueType}
	_type, ok := converter[strings.ToLower(s)]
	return _type, ok
}

package models

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/nferruzzi/gormGIS"
	"strings"
)

const FirebaseContextKey string = "firebase_context_key"

type ProfileType string

const (
	ProducerType  ProfileType = "producer"
	PerformerType             = "performer"
	VenueType                 = "venue"
)

func (t *ProfileType) UnmarshallJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	switch s {
	case "producer":
		*t = ProducerType
	case "performer":
		*t = PerformerType
	case "venue":
		*t = VenueType
	default:
		return fmt.Errorf("unrecognized profile type %s. Allowed: producer, performer, venue", s)
	}
	return nil
}

type UserID struct {
	Model
	FirebaseId  string
	Permissions Permission
	ProfileId   uuid.UUID `json:"-"`
}

type Profile struct {
	Model
	Name        string `json:"name,nonempty" gorm:"notnull"`
	ProfileType `json:"type"`
	Location    *gormGIS.GeoPoint
	UserIDs     []UserID `json:"-" gorm:"foreignKey:ProfileId"`
}

func ParseProfile(s string) (ProfileType, bool) {
	converter := map[string]ProfileType{"producer": ProducerType, "performer": PerformerType, "venue": VenueType}
	_type, ok := converter[strings.ToLower(s)]
	return _type, ok
}

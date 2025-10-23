package model

import (
	"time"

	"github.com/raymondsugiarto/coffee-api/pkg/model/concern"
)

type Notification struct {
	concern.CommonWithIDs
	OrganizationID string
	Organization   *Organization
	UserID         string
	User           *User
	RefModule      string
	RefTable       string
	RefID          string
	RefCode        string
	Description    string
	NotifyAt       time.Time
	ReadAt         time.Time
}

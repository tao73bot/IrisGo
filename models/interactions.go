package models

import (
	"github.com/google/uuid"
)

type Interaction struct {
	InteractionID uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	LeadID        uuid.UUID `gorm:"type:uuid;not null"`
	// CustomerID    uuid.UUID `gorm:"type:uuid;not null"`
	UserID        uuid.UUID `gorm:"type:uuid;not null"`
	Type          string    `gorm:"type:VARCHAR(255);check:Type IN ('call', 'email', 'meeting', 'lunch', 'other');default:'call'"`
	Notes         string    `gorm:"type:text"`
}
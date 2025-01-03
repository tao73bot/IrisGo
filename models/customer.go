package models

import (
	"github.com/google/uuid"
)

type Customer struct {
	CustomerID   uuid.UUID     `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	LeadID       uuid.UUID     `gorm:"type:uuid;not null"`
	UserID       uuid.UUID     `gorm:"type:uuid;not null"`
	Address      string        `gorm:"type:varchar(255);not null"`
	CompanyName  string        `gorm:"type:varchar(255);not null"`
	// Lead         Lead          `gorm:"foreignKey:LeadID;references:LeadID"`
	// User         User          `gorm:"foreignKey:UserID;references:UserID"`
}
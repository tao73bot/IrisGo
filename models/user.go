package models

import (
	"github.com/google/uuid"
)

type User struct {
	UserID    uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name      string     `gorm:"type:varchar(255);not null"`
	Email     string     `gorm:"type:varchar(255);unique;not null"`
	Password  string     `gorm:"type:varchar(255);not null"`
	Role      string     `gorm:"type:VARCHAR(255);check:role IN('admin', 'user');default:'user'"`
	Customer  []Customer `gorm:"foreignKey:UserID;references:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Lead      []Lead     `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
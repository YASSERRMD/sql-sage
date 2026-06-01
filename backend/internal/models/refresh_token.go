package models

import "github.com/google/uuid"

type RefreshToken struct {
	Base
	UserID    uuid.UUID `gorm:"type:uuid;not null;index" json:"-"`
	TokenHash string    `gorm:"type:text;uniqueIndex;not null" json:"-"`
	ExpiresAt int64     `gorm:"not null" json:"expiresAt"`
	Revoked   bool      `gorm:"not null;default:false;index" json:"revoked"`
}

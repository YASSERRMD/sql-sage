package models

import "github.com/google/uuid"

type Provider struct {
	Base
	UserID          uuid.UUID `gorm:"type:uuid;not null;index" json:"-"`
	Name            string    `gorm:"type:text;not null" json:"name"`
	BaseURL         string    `gorm:"type:text;not null" json:"baseUrl"`
	APIKeyEncrypted string    `gorm:"type:text;not null" json:"-"`
	APIKeyPreview   string    `gorm:"type:text" json:"apiKeyPreview"`
	ModelName       string    `gorm:"type:text;not null" json:"modelName"`
	Temperature     float64   `gorm:"not null;default:0.2" json:"temperature"`
	MaxTokens       int       `gorm:"not null;default:2048" json:"maxTokens"`
	IsDefault       bool      `gorm:"not null;default:false;index" json:"isDefault"`
}

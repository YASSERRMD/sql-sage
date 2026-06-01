package models

import (
	"github.com/google/uuid"
)

type ObjectType string

const (
	ObjectProcedure ObjectType = "procedure"
	ObjectFunction  ObjectType = "function"
	ObjectPackage   ObjectType = "package"
	ObjectTrigger   ObjectType = "trigger"
	ObjectView      ObjectType = "view"
	ObjectScript    ObjectType = "sql_script"
	ObjectUnknown   ObjectType = "unknown"
)

func (o ObjectType) Valid() bool {
	switch o {
	case ObjectProcedure, ObjectFunction, ObjectPackage,
		ObjectTrigger, ObjectView, ObjectScript, ObjectUnknown:
		return true
	}
	return false
}

type RiskLevel string

const (
	RiskLow      RiskLevel = "low"
	RiskMedium   RiskLevel = "medium"
	RiskHigh     RiskLevel = "high"
	RiskCritical RiskLevel = "critical"
)

type Analysis struct {
	Base
	UserID     uuid.UUID  `gorm:"type:uuid;not null;index" json:"-"`
	ProviderID uuid.UUID  `gorm:"type:uuid;not null;index" json:"providerId"`
	ObjectName string     `gorm:"type:text;not null;index" json:"objectName"`
	ObjectType ObjectType `gorm:"type:text;not null" json:"objectType"`
	SourceCode string     `gorm:"type:text;not null" json:"sourceCode"`
	Summary    string     `gorm:"type:text" json:"summary"`
	RiskScore  string     `gorm:"type:text;index" json:"riskScore"`
	Result     string     `gorm:"type:text;not null" json:"result"`
	TokensUsed int        `gorm:"not null;default:0" json:"tokensUsed"`
}

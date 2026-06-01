package models

const (
	RoleAdmin = "admin"
	RoleUser  = "user"
)

type User struct {
	Base
	Email        string `gorm:"type:text;uniqueIndex;not null" json:"email"`
	Name         string `gorm:"type:text;not null" json:"name"`
	PasswordHash string `gorm:"type:text;not null" json:"-"`
	Role         string `gorm:"type:text;not null;default:user" json:"role"`
	IsActive     bool   `gorm:"not null;default:true" json:"isActive"`
}

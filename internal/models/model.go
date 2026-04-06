package models

import "gorm.io/gorm"

type Users struct {
	gorm.Model
	Name     string `gorm:"not null"`
	Email    string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
	UserId   uint   `gorm:"unique;not null"`
}
type TempUser struct {
	Name     string
	Email    string
	Password string
	UserId   uint
}
type Session struct {
	SessionId     string
	TemporaryCode uint
}
